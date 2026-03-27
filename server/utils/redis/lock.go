package redis

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	DefaultLockExpiration = 30 * time.Second // 默认锁过期时间
	DefaultRetryInterval  = 100 * time.Millisecond
	DefaultMaxRetries     = 50               // 最大重试次数
	WatchdogInterval      = 10 * time.Second // 看门狗续期间隔（锁过期时间的1/3）
)

// RedisLock 分布式锁
type RedisLock struct {
	key        string
	value      string
	expiration time.Duration

	// 看门狗相关
	stopCh   chan struct{}
	stopOnce sync.Once
	isLocked bool
	mu       sync.Mutex
}

// NewRedisLock 创建一个新的 Redis 锁实例
func NewRedisLock(key string, expiration time.Duration) *RedisLock {
	if expiration == 0 {
		expiration = DefaultLockExpiration
	}
	return &RedisLock{
		key:        key,
		value:      uuid.New().String(), // 使用 UUID 保证唯一性
		expiration: expiration,
		stopCh:     make(chan struct{}),
	}
}

// Lock 尝试获取锁
func (l *RedisLock) Lock(ctx context.Context) (bool, error) {
	success, err := global.GVA_REDIS.SetNX(ctx, l.key, l.value, l.expiration).Result()
	if err != nil {
		return false, errors.Wrap(err, "failed to acquire redis lock")
	}
	if success {
		l.mu.Lock()
		l.isLocked = true
		l.mu.Unlock()
		// 启动看门狗
		go l.startWatchdog(ctx)
	}
	return success, nil
}

// Unlock 释放锁
func (l *RedisLock) Unlock(ctx context.Context) error {
	// 先停止看门狗
	l.stopWatchdog()

	l.mu.Lock()
	if !l.isLocked {
		l.mu.Unlock()
		return nil
	}
	l.isLocked = false
	l.mu.Unlock()

	// 使用 Lua 脚本保证原子性：只有 value 匹配时才删除
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	result, err := global.GVA_REDIS.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return errors.Wrap(err, "failed to release redis lock")
	}
	if result.(int64) == 0 {
		logx.Error("lock already expired or released by another process", logx.Field("key", l.key))
	}
	return nil
}

// startWatchdog 启动看门狗，定期续期锁
func (l *RedisLock) startWatchdog(ctx context.Context) {
	ticker := time.NewTicker(WatchdogInterval)
	defer ticker.Stop()

	for {
		select {
		case <-l.stopCh:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.mu.Lock()
			if !l.isLocked {
				l.mu.Unlock()
				return
			}
			l.mu.Unlock()

			// 使用 Lua 脚本原子性续期：只有 value 匹配时才续期
			script := `
				if redis.call("get", KEYS[1]) == ARGV[1] then
					return redis.call("pexpire", KEYS[1], ARGV[2])
				else
					return 0
				end
			`
			result, err := global.GVA_REDIS.Eval(context.Background(), script,
				[]string{l.key}, l.value, int64(l.expiration/time.Millisecond)).Result()
			if err != nil {
				logx.Error("watchdog renew lock failed", logx.Field("err", err), logx.Field("key", l.key))
				return
			}
			if result.(int64) == 0 {
				logx.Error("watchdog: lock expired or stolen", logx.Field("key", l.key))
				l.mu.Lock()
				l.isLocked = false
				l.mu.Unlock()
				return
			}

			logx.Info("watchdog renewed lock", logx.Field("key", l.key))
		}
	}
}

// stopWatchdog 停止看门狗
func (l *RedisLock) stopWatchdog() {
	l.stopOnce.Do(func() {
		close(l.stopCh)
	})
}

// LockWithRetry 尝试获取锁，带有重试机制
func (l *RedisLock) LockWithRetry(ctx context.Context, maxRetries int, retryInterval time.Duration) (bool, error) {
	if maxRetries <= 0 {
		maxRetries = DefaultMaxRetries
	}
	if retryInterval == 0 {
		retryInterval = DefaultRetryInterval
	}

	for i := 0; i < maxRetries; i++ {
		success, err := l.Lock(ctx)
		if err != nil {
			return false, err
		}
		if success {
			return true, nil
		}

		// 检查 context 是否已取消
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(retryInterval):
			continue
		}
	}
	return false, nil
}

// AcquireProductLock 获取产品库存锁
// key: product_inventory_lock:{productId}
func AcquireProductLock(ctx context.Context, productId uint) (*RedisLock, error) {
	key := fmt.Sprintf("product_inventory_lock:%d", productId)
	lock := NewRedisLock(key, DefaultLockExpiration)

	success, err := lock.LockWithRetry(ctx, DefaultMaxRetries, DefaultRetryInterval)
	if err != nil {
		logx.Error("获取分布式锁出错", logx.Field("err", err))
		return nil, err
	}
	if !success {
		return nil, errors.New("系统繁忙，请稍后重试")
	}
	return lock, nil
}

// AcquireResourceLock 通用资源锁
// resourceType: 资源类型（如 "product", "volume", "training" 等）
// resourceId: 资源ID
func AcquireResourceLock(ctx context.Context, resourceType string, resourceId uint) (*RedisLock, error) {
	key := fmt.Sprintf("resource_lock:%s:%d", resourceType, resourceId)
	lock := NewRedisLock(key, DefaultLockExpiration)

	success, err := lock.LockWithRetry(ctx, DefaultMaxRetries, DefaultRetryInterval)
	if err != nil {
		logx.Error("获取分布式锁出错", logx.Field("err", err), logx.Field("resourceType", resourceType), logx.Field("resourceId", resourceId))
		return nil, err
	}
	if !success {
		return nil, errors.New("系统繁忙，请稍后重试")
	}
	return lock, nil
}
