package system

import (
	"context"
	"gin-vue-admin/global"
	"gin-vue-admin/model/system"
	"gin-vue-admin/utils"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type JwtService struct{}

var JwtServiceApp = new(JwtService)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList model.JwtBlacklist
//@return: err error

func (jwtService *JwtService) JsonInBlacklist(jwtList system.JwtBlacklist) (err error) {
	err = global.GVA_DB.Create(&jwtList).Error
	if err != nil {
		return
	}
	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})
	return
}

// GetRedisJWT 从redis取当前用户的一个活跃jwt
func (jwtService *JwtService) GetRedisJWT(userName string) (redisJWT string, err error) {
	ctx := context.Background()
	key := utils.GetUserSessionsKey(userName)
	// 每次查询前，先异步清理已过期的 Token
	go global.GVA_REDIS.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(time.Now().Unix(), 10))

	// 获取最新的一个有效 Token
	results, err := global.GVA_REDIS.ZRevRange(ctx, key, 0, 0).Result()
	if err != nil || len(results) == 0 {
		return "", err
	}
	return results[0], nil
}

// GetRedisJWTs 获取该用户所有活跃且未过期的 Token
func (jwtService *JwtService) GetRedisJWTs(userName string) (tokens []string, err error) {
	ctx := context.Background()
	key := utils.GetUserSessionsKey(userName)
	// 清理已过期的 Token
	now := strconv.FormatInt(time.Now().Unix(), 10)
	global.GVA_REDIS.ZRemRangeByScore(ctx, key, "0", now)

	// 获取所有剩余的 Token
	tokens, err = global.GVA_REDIS.ZRange(ctx, key, 0, -1).Result()
	return tokens, err
}

// ClearRedisJWTs 清除该用户所有的活跃会话
func (jwtService *JwtService) ClearRedisJWTs(userName string) (err error) {
	ctx := context.Background()
	key := utils.GetUserSessionsKey(userName)
	return global.GVA_REDIS.Del(ctx, key).Err()
}

func LoadAll() {
	var data []string
	err := global.GVA_DB.Model(&system.JwtBlacklist{}).Select("jwt").Find(&data).Error
	if err != nil {
		logx.Error("加载数据库jwt黑名单失败!", err)
		return
	}
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	} // jwt黑名单 加入 BlackCache 中
}
