package sshkey

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"gin-vue-admin/global"
	secretModel "gin-vue-admin/model/secret"
	secretReq "gin-vue-admin/model/secret/request"
	secretResp "gin-vue-admin/model/secret/response"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SSHKeyService struct{}

// CreateSSHKey 创建SSH密钥
func (s *SSHKeyService) CreateSSHKey(ctx context.Context, req *secretReq.AddSSHKeyReq, userId uint) (*secretResp.AddSSHKeyResp, error) {
	// 验证公钥格式
	if err := validatePublicKey(req.PublicKey); err != nil {
		return nil, err
	}

	// 生成指纹
	fingerprint, err := generateFingerprint(req.PublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "生成密钥指纹失败")
	}

	// 检查密钥是否已存在
	var existingKey secretModel.SSHKey
	if err := global.GVA_DB.Where("user_id = ? AND fingerprint = ?", userId, fingerprint).First(&existingKey).Error; err == nil {
		return nil, errors.New("该密钥已存在")
	}

	// 创建密钥记录
	sshKey := secretModel.SSHKey{
		Name:        req.Name,
		UserId:      userId,
		PublicKey:   req.PublicKey,
		Fingerprint: fingerprint,
		IsDefault:   false,
	}

	if err := global.GVA_DB.Create(&sshKey).Error; err != nil {
		logx.Error("创建SSH密钥失败", err)
		return nil, errors.Wrap(err, "创建SSH密钥失败")
	}

	return &secretResp.AddSSHKeyResp{
		ID:          sshKey.ID,
		Fingerprint: fingerprint,
	}, nil
}

// GetSSHKeyList 获取SSH密钥列表
func (s *SSHKeyService) GetSSHKeyList(ctx context.Context, req *secretReq.GetSSHKeyListReq, userId uint) (*secretResp.GetSSHKeyListResp, error) {
	var sshKeys []secretModel.SSHKey
	var total int64

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	db := global.GVA_DB.Model(&secretModel.SSHKey{}).Where("user_id = ?", userId)
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		logx.Error("查询SSH密钥总数失败", err)
		return nil, errors.Wrap(err, "查询SSH密钥总数失败")
	}

	if err := db.Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).Order("is_default DESC, id DESC").Find(&sshKeys).Error; err != nil {
		logx.Error("查询SSH密钥列表失败", err)
		return nil, errors.Wrap(err, "查询SSH密钥列表失败")
	}

	// 转换为响应格式
	var list []secretResp.SSHKeyItem
	for _, key := range sshKeys {
		list = append(list, secretResp.SSHKeyItem{
			ID:          key.ID,
			Name:        key.Name,
			Fingerprint: key.Fingerprint,
			IsDefault:   key.IsDefault,
			CreatedAt:   key.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &secretResp.GetSSHKeyListResp{
		List:  list,
		Total: total,
	}, nil
}

// DeleteSSHKey 删除SSH密钥
func (s *SSHKeyService) DeleteSSHKey(ctx context.Context, req *secretReq.DeleteSSHKeyReq, userId uint) error {
	// 检查密钥是否存在且属于当前用户
	var sshKey secretModel.SSHKey
	if err := global.GVA_DB.Where("id = ? AND user_id = ?", req.ID, userId).First(&sshKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("SSH密钥不存在")
		}
		return errors.Wrap(err, "查询SSH密钥失败")
	}

	// 删除密钥
	if err := global.GVA_DB.Delete(&sshKey).Error; err != nil {
		logx.Error("删除SSH密钥失败", err)
		return errors.Wrap(err, "删除SSH密钥失败")
	}

	return nil
}

// SetDefaultSSHKey 设置默认密钥
func (s *SSHKeyService) SetDefaultSSHKey(ctx context.Context, req *secretReq.SetDefaultSSHKeyReq, userId uint) error {
	// 开启事务
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 检查密钥是否存在且属于当前用户
		var sshKey secretModel.SSHKey
		if err := tx.Where("id = ? AND user_id = ?", req.ID, userId).First(&sshKey).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("SSH密钥不存在")
			}
			return errors.Wrap(err, "查询SSH密钥失败")
		}

		// 取消当前用户的所有默认密钥
		if err := tx.Model(&secretModel.SSHKey{}).Where("user_id = ?", userId).Update("is_default", false).Error; err != nil {
			return errors.Wrap(err, "取消默认密钥失败")
		}

		// 设置新的默认密钥
		if err := tx.Model(&sshKey).Update("is_default", true).Error; err != nil {
			return errors.Wrap(err, "设置默认密钥失败")
		}

		return nil
	})
}

// validatePublicKey 验证公钥格式
func validatePublicKey(publicKey string) error {
	publicKey = strings.TrimSpace(publicKey)
	if publicKey == "" {
		return errors.New("公钥不能为空")
	}

	// 检查是否是有效的SSH公钥格式
	parts := strings.Fields(publicKey)
	if len(parts) < 2 {
		return errors.New("公钥格式不正确")
	}

	// 检查密钥类型
	keyType := parts[0]
	validTypes := []string{"ssh-rsa", "ssh-dss", "ssh-ed25519", "ecdsa-sha2-nistp256", "ecdsa-sha2-nistp384", "ecdsa-sha2-nistp521"}
	isValid := false
	for _, t := range validTypes {
		if keyType == t {
			isValid = true
			break
		}
	}
	if !isValid {
		return errors.New("不支持的密钥类型")
	}

	// 验证base64编码
	if _, err := base64.StdEncoding.DecodeString(parts[1]); err != nil {
		return errors.New("公钥格式不正确")
	}

	return nil
}

// generateFingerprint 生成SSH密钥指纹
func generateFingerprint(publicKey string) (string, error) {
	parts := strings.Fields(publicKey)
	if len(parts) < 2 {
		return "", errors.New("公钥格式不正确")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.Wrap(err, "解码公钥失败")
	}

	hash := md5.Sum(decoded)
	fingerprint := fmt.Sprintf("%x", hash)

	// 格式化为标准指纹格式 (xx:xx:xx:...)
	var formatted strings.Builder
	for i := 0; i < len(fingerprint); i += 2 {
		if i > 0 {
			formatted.WriteString(":")
		}
		formatted.WriteString(fingerprint[i : i+2])
	}

	return formatted.String(), nil
}
