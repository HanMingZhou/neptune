package inference

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"gin-vue-admin/global"
	inferenceModel "gin-vue-admin/model/inference"
	inferenceReq "gin-vue-admin/model/inference/request"
	"gin-vue-admin/model/inference/response"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// InferenceApiKeyService API Key 服务
type InferenceApiKeyService struct{}

var InferenceApiKeyServiceApp = new(InferenceApiKeyService)

// CreateApiKey 创建 API Key
func (s *InferenceApiKeyService) CreateApiKey(ctx context.Context, req *inferenceReq.CreateApiKeyReq) (*response.CreateApiKeyResp, error) {
	// 验证服务是否存在
	var service inferenceModel.Inference
	if err := global.GVA_DB.Where("id = ?", req.ServiceId).First(&service).Error; err != nil {
		return nil, errors.Wrap(err, "服务不存在")
	}

	// 验证用户权限 (只能为自己的服务创建 Key)
	if service.UserId != req.UserId {
		return nil, errors.New("无权为此服务创建 API Key")
	}

	// 生成 API Key
	apiKey, err := generateApiKey()
	if err != nil {
		return nil, errors.Wrap(err, "生成API Key失败")
	}

	// 设置过期时间
	var expiredAt *time.Time
	if req.ExpireDays > 0 {
		t := time.Now().AddDate(0, 0, req.ExpireDays)
		expiredAt = &t
	}

	// 设置默认权限范围
	scopes := req.Scopes
	if scopes == "" {
		scopes = "read,write"
	}

	// 创建记录
	key := &inferenceModel.InferenceApiKey{
		ServiceId:   req.ServiceId,
		ApiKey:      apiKey,
		Name:        req.Name,
		Description: req.Description,
		Status:      inferenceModel.ApiKeyStatusActive,
		Scopes:      scopes,
		RateLimit:   req.RateLimit,
		ExpiredAt:   expiredAt,
		UserId:      req.UserId,
	}

	if err := global.GVA_DB.Create(key).Error; err != nil {
		return nil, errors.Wrap(err, "创建API Key失败")
	}

	return &response.CreateApiKeyResp{
		ID:     key.ID,
		ApiKey: apiKey, // 仅创建时返回完整 Key
		Name:   key.Name,
	}, nil
}

// ListApiKeys 获取 API Key 列表
func (s *InferenceApiKeyService) ListApiKeys(ctx context.Context, req *inferenceReq.ListApiKeysReq) (*response.ListApiKeysResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	db := global.GVA_DB.Model(&inferenceModel.InferenceApiKey{}).Where("service_id = ?", req.ServiceId)

	var total int64
	db.Count(&total)

	var keys []inferenceModel.InferenceApiKey
	if err := db.Order("created_at DESC").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&keys).Error; err != nil {
		return nil, err
	}

	list := make([]response.ApiKeyItem, len(keys))
	for i, k := range keys {
		list[i] = response.ApiKeyItem{
			ID:          k.ID,
			Name:        k.Name,
			ApiKey:      maskApiKey(k.ApiKey), // 脱敏显示
			Description: k.Description,
			Status:      k.Status,
			Scopes:      k.Scopes,
			RateLimit:   k.RateLimit,
			LastUsedAt:  k.LastUsedAt,
			ExpiredAt:   k.ExpiredAt,
			CreatedAt:   k.CreatedAt,
		}
	}

	return &response.ListApiKeysResp{
		Total: total,
		List:  list,
	}, nil
}

// DeleteApiKey 删除 API Key
func (s *InferenceApiKeyService) DeleteApiKey(ctx context.Context, req *inferenceReq.DeleteApiKeyReq) error {
	return global.GVA_DB.Delete(&inferenceModel.InferenceApiKey{}, req.ID).Error
}

// generateApiKey 生成 API Key
func generateApiKey() (string, error) {
	bytes := make([]byte, 24)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "sk-" + hex.EncodeToString(bytes), nil
}

// maskApiKey 脱敏显示 API Key
func maskApiKey(key string) string {
	if len(key) <= 10 {
		return key
	}
	return key[:6] + strings.Repeat("*", len(key)-10) + key[len(key)-4:]
}
