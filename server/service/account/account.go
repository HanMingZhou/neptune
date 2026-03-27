package account

import (
	"bytes"
	"context"
	"encoding/base64"
	"gin-vue-admin/global"
	"gin-vue-admin/model/account"
	accountReq "gin-vue-admin/model/account/request"
	accountResp "gin-vue-admin/model/account/response"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/service/system"
	"gin-vue-admin/utils"
	"image/png"
	"time"

	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
)

type AccountService struct{}

func (s *AccountService) GetAccessLogList(ctx context.Context, req accountReq.GetAccessLogListReq, userId uint) (resp accountResp.GetAccessLogListResp, err error) {
	var logs []account.AccountAccessLog
	var total int64
	db := global.GVA_DB.Model(&account.AccountAccessLog{}).Where("user_id = ?", userId)

	if req.Ip != "" {
		db = db.Where("ip LIKE ?", "%"+req.Ip+"%")
	}
	if req.Device != "" {
		db = db.Where("device LIKE ?", "%"+req.Device+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	err = db.Count(&total).Error
	if err != nil {
		return resp, err
	}

	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	err = db.Limit(limit).Offset(offset).Order("created_at desc").Find(&logs).Error
	if err != nil {
		return resp, err
	}

	resp.Total = total
	resp.List = make([]accountResp.AccessLogInfo, 0, len(logs))
	for _, l := range logs {
		resp.List = append(resp.List, accountResp.AccessLogInfo{
			ID:       l.ID,
			Time:     l.CreatedAt.Format("2006-01-02 15:04:05"),
			Ip:       l.Ip,
			Location: l.Location,
			Device:   l.Device,
			Os:       l.Os,
			Browser:  l.Browser,
			Method:   l.Method,
			Status:   l.Status,
			Reason:   l.Reason,
			LogType:  l.LogType,
		})
	}

	return resp, nil
}

func (s *AccountService) GetActiveSessionList(ctx context.Context, req accountReq.GetAccessLogListReq, userId uint) (resp accountResp.GetAccessLogListResp, err error) {
	var user sysModel.SysUser
	err = global.GVA_DB.First(&user, userId).Error
	if err != nil {
		return resp, err
	}

	if global.GVA_REDIS == nil {
		return resp, nil
	}

	tokens, err := system.JwtServiceApp.GetRedisJWTs(user.Username)
	if err != nil {
		return resp, err
	}

	activeCount := len(tokens)
	if activeCount == 0 {
		resp.Total = 0
		resp.List = []accountResp.AccessLogInfo{}
		return resp, nil
	}

	var logs []account.AccountAccessLog
	err = global.GVA_DB.Where("user_id = ? AND log_type = ? AND status = ?", userId, account.LogTypeLogin, account.LogStatusSuccess).
		Order("created_at desc").
		Limit(activeCount).
		Find(&logs).Error

	if err != nil {
		return resp, err
	}

	resp.Total = int64(activeCount)
	resp.List = make([]accountResp.AccessLogInfo, 0, activeCount)
	for i, l := range logs {
		resp.List = append(resp.List, accountResp.AccessLogInfo{
			ID:       l.ID,
			Time:     l.CreatedAt.Format("2006-01-02 15:04:05"),
			Ip:       l.Ip,
			Location: l.Location,
			Device:   l.Device,
			Os:       l.Os,
			Browser:  l.Browser,
			Method:   l.Method,
			Status:   account.LogStatusSuccess,
			Reason:   "Active Session",
			LogType:  l.LogType,
		})
		if i+1 >= activeCount {
			break
		}
	}

	return resp, nil
}

func (s *AccountService) GetSecurityStatus(ctx context.Context, userId uint) (resp accountResp.SecurityStatusResp, err error) {
	var user sysModel.SysUser
	err = global.GVA_DB.First(&user, userId).Error
	if err != nil {
		return resp, err
	}

	resp.Phone = user.Phone
	if resp.Phone != "" {
		resp.PhoneStatus = account.StatusBound
	} else {
		resp.PhoneStatus = account.StatusUnset
	}

	resp.Email = user.Email
	if resp.Email != "" {
		resp.EmailStatus = account.StatusBound
	} else {
		resp.EmailStatus = account.StatusUnset
	}

	var security account.AccountSecurity
	err = global.GVA_DB.FirstOrCreate(&security, account.AccountSecurity{UserId: userId}).Error
	if err != nil {
		return resp, err
	}

	resp.MfaEnabled = security.MfaEnabled
	if resp.MfaEnabled {
		resp.MfaStatus = account.StatusMfaEnabled
	} else {
		resp.MfaStatus = account.StatusMfaDisabled
	}

	if security.GithubId != "" {
		resp.GithubBound = true
		resp.GithubUsername = security.GithubUsername
		resp.GithubStatus = account.StatusLinked
	} else {
		resp.GithubStatus = account.StatusUnlinked
	}

	if security.AccessKeyId != "" {
		resp.AccessKeyId = security.AccessKeyId
		resp.AccessKeyStatus = account.StatusKeyGenerated
	} else {
		resp.AccessKeyStatus = account.StatusKeyNotGenerated
	}

	score := 60
	if resp.Phone != "" {
		score += 10
	}
	if resp.Email != "" {
		score += 10
	}
	if resp.MfaEnabled {
		score += 20
	}
	resp.SecurityScore = score

	var lastLog account.AccountAccessLog
	if err := global.GVA_DB.Where("user_id = ? AND log_type = ?", userId, account.LogTypeLogin).Order("created_at desc").First(&lastLog).Error; err == nil {
		resp.LastLoginTime = lastLog.CreatedAt.Format("2006-01-02 15:04:05")
		resp.LastLoginIp = lastLog.Ip
	} else {
		resp.LastLoginTime = time.Now().Format("2006-01-02 15:04:05")
		resp.LastLoginIp = "Unknown"
	}

	return resp, nil
}

func (s *AccountService) UpdatePassword(userId uint, req accountReq.UpdatePasswordReq) error {
	var user sysModel.SysUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return err
	}
	if !utils.BcryptCheck(req.OldPassword, user.Password) {
		return errors.New("原密码不正确")
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.GVA_DB.Save(&user).Error
}

func (s *AccountService) BindAccount(userId uint, req accountReq.BindAccountReq) error {
	// TODO: 校验验证码 req.Code
	switch req.Type {
	case account.BindTypePhone:
		return global.GVA_DB.Model(&sysModel.SysUser{}).Where("id = ?", userId).Update("phone", req.Value).Error
	case account.BindTypeEmail:
		return global.GVA_DB.Model(&sysModel.SysUser{}).Where("id = ?", userId).Update("email", req.Value).Error
	default:
		return errors.New("未知绑定类型")
	}
}

// MfaSetup 生成TOTP密钥和二维码（尚未激活，需用户扫码后调用MfaActivate确认）
func (s *AccountService) MfaSetup(userId uint) (resp accountResp.MfaSecretResp, err error) {
	var user sysModel.SysUser
	if err = global.GVA_DB.First(&user, userId).Error; err != nil {
		return resp, errors.New("用户不存在")
	}

	// 生成TOTP密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "NeptunePlatform",
		AccountName: user.Username,
	})
	if err != nil {
		return resp, errors.Wrap(err, "生成MFA密钥失败")
	}

	// 生成二维码图片并转为base64
	img, err := key.Image(200, 200)
	if err != nil {
		return resp, errors.Wrap(err, "生成二维码失败")
	}
	var buf bytes.Buffer
	if err = png.Encode(&buf, img); err != nil {
		return resp, errors.Wrap(err, "编码二维码失败")
	}

	// 将密钥临时存入数据库（MfaEnabled仍为false，等待用户确认激活）
	var security account.AccountSecurity
	if err = global.GVA_DB.FirstOrCreate(&security, account.AccountSecurity{UserId: userId}).Error; err != nil {
		return resp, err
	}
	if err = global.GVA_DB.Model(&account.AccountSecurity{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"mfa_secret":  key.Secret(),
		"mfa_enabled": false,
	}).Error; err != nil {
		return resp, err
	}

	resp.Secret = key.Secret()
	resp.Qr = "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	return resp, nil
}

// MfaActivate 用户扫码后输入验证码，验证通过则激活MFA
func (s *AccountService) MfaActivate(userId uint, code string) error {
	var security account.AccountSecurity
	if err := global.GVA_DB.Where("user_id = ?", userId).First(&security).Error; err != nil {
		return errors.New("请先初始化MFA")
	}
	if security.MfaSecret == "" {
		return errors.New("请先初始化MFA")
	}
	if security.MfaEnabled {
		return errors.New("MFA已开启，无需重复激活")
	}

	// 验证TOTP码
	if !totp.Validate(code, security.MfaSecret) {
		return errors.New("验证码错误，请重试")
	}

	return global.GVA_DB.Model(&account.AccountSecurity{}).Where("user_id = ?", userId).Update("mfa_enabled", true).Error
}

// ToggleMfa 关闭MFA（需验证当前TOTP码）
func (s *AccountService) ToggleMfa(userId uint, req accountReq.ToggleMfaReq) error {
	if req.Enabled {
		return errors.New("开启MFA请使用MFA初始化流程")
	}

	var security account.AccountSecurity
	if err := global.GVA_DB.Where("user_id = ?", userId).First(&security).Error; err != nil {
		return errors.New("MFA记录不存在")
	}
	if !security.MfaEnabled {
		return errors.New("MFA尚未开启")
	}

	// 验证TOTP码
	if !totp.Validate(req.Code, security.MfaSecret) {
		return errors.New("验证码错误")
	}

	return global.GVA_DB.Model(&account.AccountSecurity{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"mfa_enabled": false,
		"mfa_secret":  "",
	}).Error
}

// VerifyMfaCode 登录时验证TOTP码
func (s *AccountService) VerifyMfaCode(userId uint, code string) error {
	var security account.AccountSecurity
	if err := global.GVA_DB.Where("user_id = ?", userId).First(&security).Error; err != nil {
		return errors.New("MFA记录不存在")
	}
	if !security.MfaEnabled || security.MfaSecret == "" {
		return errors.New("用户未开启MFA")
	}
	if !totp.Validate(code, security.MfaSecret) {
		return errors.New("MFA验证码错误")
	}
	return nil
}

// IsMfaEnabled 检查用户是否开启了MFA
func (s *AccountService) IsMfaEnabled(userId uint) bool {
	var security account.AccountSecurity
	if err := global.GVA_DB.Where("user_id = ?", userId).First(&security).Error; err != nil {
		return false
	}
	return security.MfaEnabled
}

func (s *AccountService) GenerateAccessKey(userId uint) (resp accountResp.GenerateAccessKeyResp, err error) {
	resp.AccessKeyId = utils.RandomString(16)
	resp.AccessKeySecret = utils.RandomString(32)

	var security account.AccountSecurity
	err = global.GVA_DB.FirstOrCreate(&security, account.AccountSecurity{UserId: userId}).Error
	if err != nil {
		return resp, err
	}

	err = global.GVA_DB.Model(&account.AccountSecurity{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"access_key_id":     resp.AccessKeyId,
		"access_key_secret": resp.AccessKeySecret,
	}).Error
	return resp, err
}

func (s *AccountService) KillSession(ctx context.Context, req accountReq.KillSessionReq, userId uint) error {
	var user sysModel.SysUser
	if err := global.GVA_DB.First(&user, userId).Error; err != nil {
		return err
	}

	tokens, err := system.JwtServiceApp.GetRedisJWTs(user.Username)
	if err != nil || len(tokens) == 0 {
		return err
	}

	var logs []account.AccountAccessLog
	global.GVA_DB.Where("user_id = ? AND log_type = ? AND status = ?", userId, account.LogTypeLogin, account.LogStatusSuccess).
		Order("created_at desc").
		Limit(len(tokens)).
		Find(&logs)

	targetIdx := -1
	for i, l := range logs {
		if l.ID == req.LogId {
			targetIdx = i
			break
		}
	}

	if targetIdx == -1 {
		return errors.New("未找到对应会话")
	}

	tokenToKill := tokens[len(tokens)-1-targetIdx]

	_ = system.JwtServiceApp.JsonInBlacklist(sysModel.JwtBlacklist{Jwt: tokenToKill})
	return utils.DelRedisJWT(user.Username, tokenToKill)
}

func (s *AccountService) CreateAccessLog(log account.AccountAccessLog) error {
	return global.GVA_DB.Create(&log).Error
}
