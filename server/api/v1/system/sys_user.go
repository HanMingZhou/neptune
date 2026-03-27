package system

import (
	"context"
	cryptoRand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"gin-vue-admin/global"
	accountModel "gin-vue-admin/model/account"
	accountReq "gin-vue-admin/model/account/request"
	"gin-vue-admin/model/common"
	"gin-vue-admin/model/common/request"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/system"
	systemReq "gin-vue-admin/model/system/request"
	systemRes "gin-vue-admin/model/system/response"
	"gin-vue-admin/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Login
// @Tags     Base
// @Summary  用户登录
// @Produce   application/json
// @Param    data  body      systemReq.Login                                             true  "用户名, 密码, 验证码"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "返回包括用户信息,token,过期时间"
// @Router   /api/v1/base/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var l systemReq.Login
	err := c.ShouldBindJSON(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.LoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	key := c.ClientIP()
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha
	openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool = openCaptcha == 0 || openCaptcha < interfaceToInt(v)
	if oc && (l.Captcha == "" || l.CaptchaId == "" || !store.Verify(l.CaptchaId, l.Captcha, true)) {
		global.BlackCache.Increment(key, 1)
		response.FailWithMessage("验证码错误", c)
		return
	}

	user, err := userService.Login(l.Username, l.Password)
	if err != nil {
		ua := c.Request.UserAgent()
		go func() {
			_ = accountService.CreateAccessLog(accountModel.AccountAccessLog{
				UserId:    0,
				Ip:        c.ClientIP(),
				Location:  "Unknown",
				Device:    utils.GetBrowser(ua) + " / " + utils.GetOS(ua),
				Browser:   utils.GetBrowser(ua),
				Os:        utils.GetOS(ua),
				Method:    accountModel.LogMethodPassword,
				Status:    accountModel.LogStatusFailed,
				Reason:    err.Error(),
				LogType:   accountModel.LogTypeLogin,
				LoginTime: time.Now(),
			})
		}()
		global.BlackCache.Increment(key, 1)
		utils.HandleError(c, err, "登陆失败! 用户名不存在或者密码错误")
		return
	}

	if user.Enable != 1 {
		global.BlackCache.Increment(key, 1)
		response.FailWithMessage("用户被禁止登录", c)
		return
	}

	// 检查用户是否开启了MFA
	if accountService.IsMfaEnabled(user.ID) {
		// 生成临时MFA令牌，存入Redis，有效期5分钟
		mfaToken, err := generateMfaToken()
		if err != nil {
			utils.HandleError(c, err, "生成MFA令牌失败")
			return
		}
		// 将用户信息序列化后存入Redis
		userBytes, _ := json.Marshal(user)
		if global.GVA_REDIS != nil {
			global.GVA_REDIS.Set(context.Background(), "mfa_token:"+mfaToken, string(userBytes), 5*time.Minute)
		}
		response.OkWithDetailed(systemRes.LoginResponse{
			NeedMfa:  true,
			MfaToken: mfaToken,
		}, "需要MFA验证", c)
		return
	}

	b.TokenNext(c, user)
}

// generateMfaToken 生成安全的临时MFA令牌
func generateMfaToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := cryptoRand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// MfaLogin MFA二次验证登录
// @Tags     Base
// @Summary  MFA验证码登录
// @Produce  application/json
// @Param    data  body      accountReq.MfaLoginReq  true  "MFA临时令牌和验证码"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}
// @Router   /api/v1/base/mfa/login [post]
func (b *BaseApi) MfaLogin(c *gin.Context) {
	var req accountReq.MfaLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从Redis获取临时令牌对应的用户信息
	if global.GVA_REDIS == nil {
		response.FailWithMessage("系统未启用Redis，无法使用MFA", c)
		return
	}

	redisKey := "mfa_token:" + req.MfaToken
	userJson, err := global.GVA_REDIS.Get(context.Background(), redisKey).Result()
	if err != nil {
		response.FailWithMessage("MFA令牌无效或已过期，请重新登录", c)
		return
	}

	var user systemRes.UserResponse
	if err := json.Unmarshal([]byte(userJson), &user); err != nil {
		response.FailWithMessage("用户信息解析失败", c)
		return
	}

	// 验证TOTP码
	if err := accountService.VerifyMfaCode(user.ID, req.Code); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证通过，删除临时令牌
	global.GVA_REDIS.Del(context.Background(), redisKey)

	// 签发正式JWT
	b.TokenNext(c, user)
}

// TokenNext 登录以后签发jwt
func (b *BaseApi) TokenNext(c *gin.Context, user systemRes.UserResponse) {
	ua := c.Request.UserAgent()
	go func() {
		_ = accountService.CreateAccessLog(accountModel.AccountAccessLog{
			UserId:    user.ID,
			Ip:        c.ClientIP(),
			Location:  "Unknown",
			Device:    utils.GetBrowser(ua) + " / " + utils.GetOS(ua),
			Browser:   utils.GetBrowser(ua),
			Os:        utils.GetOS(ua),
			Method:    accountModel.LogMethodPassword,
			Status:    accountModel.LogStatusSuccess,
			LogType:   accountModel.LogTypeLogin,
			LoginTime: time.Now(),
		})
	}()

	// 转换为SysUser用于JWT签发
	u, _ := uuid.Parse(user.UUID)
	sysUser := system.SysUser{
		GVA_MODEL:   global.GVA_MODEL{ID: user.ID},
		UUID:        u,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
		Namespace:   user.Namespace,
	}

	token, claims, err := utils.LoginToken(&sysUser)
	if err != nil {
		utils.HandleError(c, err, "获取token失败")
		return
	}

	if global.GVA_CONFIG.System.UseMultipoint {
		tokens, _ := jwtService.GetRedisJWTs(user.Username)
		for _, t := range tokens {
			_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: t})
		}
		_ = jwtService.ClearRedisJWTs(user.Username)
	}

	if global.GVA_CONFIG.System.UseRedis {
		if err := utils.SetRedisJWT(token, user.Username); err != nil {
			utils.HandleError(c, err, "设置登录状态失败")
			return
		}
	}

	utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
	response.OkWithDetailed(systemRes.LoginResponse{
		User:      user,
		Token:     token,
		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
	}, "登录成功", c)
}

// Register
// @Tags     SysUser
// @Summary  用户注册账号
// @Produce   application/json
// @Param    data  body      systemReq.Register                                            true  "用户名, 昵称, 密码, 角色ID"
// @Success  200   {object}  response.Response{data=systemRes.SysUserResponse,msg=string}  "用户注册账号,返回包括用户信息"
// @Router   /api/v1/user/register [post]
func (b *BaseApi) Register(c *gin.Context) {
	var r systemReq.Register
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userResp, err := userService.Register(r)
	if err != nil {
		utils.HandleError(c, err, "注册失败")
		return
	}
	response.OkWithDetailed(systemRes.SysUserResponse{User: userResp}, "注册成功", c)
}

// ChangePassword
// @Tags      SysUser
// @Summary   用户修改密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      systemReq.ChangePasswordReq    true  "用户名, 原密码, 新密码"
// @Success   200   {object}  response.Response{msg=string}  "用户修改密码"
// @Router    /api/v1/user/password/change [post]
func (b *BaseApi) ChangePassword(c *gin.Context) {
	var req systemReq.ChangePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.ChangePasswordVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	err = userService.ChangePassword(uid, req.Password, req.NewPassword)
	if err != nil {
		utils.HandleError(c, err, "修改失败")
		return
	}
	response.OkWithMessage("修改成功", c)
}

// GetUserList
// @Tags      SysUser
// @Summary   分页获取用户列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.GetUserList                                        true  "页码, 每页大小"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router    /api/v1/user/list [post]
func (b *BaseApi) GetUserList(c *gin.Context) {
	var pageInfo systemReq.GetUserList
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := userService.GetUserInfoList(pageInfo)
	if err != nil {
		utils.HandleError(c, err, "获取失败")
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// SetUserAuthority
// @Tags      SysUser
// @Summary   更改用户权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuth          true  "用户UUID, 角色ID"
// @Success   200   {object}  response.Response{msg=string}  "设置用户权限"
// @Router    /api/v1/user/authority/set [post]
func (b *BaseApi) SetUserAuthority(c *gin.Context) {
	var sua systemReq.SetUserAuth
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if UserVerifyErr := utils.Verify(sua, utils.SetUserAuthorityVerify); UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	userID := utils.GetUserID(c)
	err = userService.SetUserAuthority(userID, sua.AuthorityId)
	if err != nil {
		utils.HandleError(c, err, "修改失败")
		return
	}
	claims := utils.GetUserInfo(c)
	claims.AuthorityId = sua.AuthorityId
	token, err := utils.NewJWT().CreateToken(*claims)
	if err != nil {
		utils.HandleError(c, err, "修改失败")
		return
	}
	c.Header("new-token", token)
	c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt.Unix(), 10))
	utils.SetToken(c, token, int(claims.ExpiresAt.Unix()-time.Now().Unix()))
	response.OkWithMessage("修改成功", c)
}

// SetUserAuthorities
// @Tags      SysUser
// @Summary   设置用户权限
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuthorities   true  "用户UUID, 角色ID"
// @Success   200   {object}  response.Response{msg=string}  "设置用户权限"
// @Router    /api/v1/user/authorities/set [post]
func (b *BaseApi) SetUserAuthorities(c *gin.Context) {
	var sua systemReq.SetUserAuthorities
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	err = userService.SetUserAuthorities(authorityID, sua.ID, sua.AuthorityIds)
	if err != nil {
		utils.HandleError(c, err, "修改失败")
		return
	}
	response.OkWithMessage("修改成功", c)
}

// DeleteUser
// @Tags      SysUser
// @Summary   删除用户
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "用户ID"
// @Success   200   {object}  response.Response{msg=string}  "删除用户"
// @Router    /api/v1/user/delete [post]
func (b *BaseApi) DeleteUser(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	jwtId := utils.GetUserID(c)
	if jwtId == uint(reqId.ID) {
		response.FailWithMessage("删除失败, 无法删除自己。", c)
		return
	}
	err = userService.DeleteUser(reqId.ID)
	if err != nil {
		utils.HandleError(c, err, "删除失败")
		return
	}
	response.OkWithMessage("删除成功", c)
}

// SetUserInfo
// @Tags      SysUser
// @Summary   设置用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /api/v1/user/info/set [post]
func (b *BaseApi) SetUserInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(user, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(user.AuthorityIds) != 0 {
		authorityID := utils.GetUserAuthorityId(c)
		err = userService.SetUserAuthorities(authorityID, user.ID, user.AuthorityIds)
		if err != nil {
			utils.HandleError(c, err, "设置失败")
			return
		}
	}
	err = userService.SetUserInfo(system.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		Enable:    user.Enable,
	})
	if err != nil {
		utils.HandleError(c, err, "设置失败")
		return
	}
	response.OkWithMessage("设置成功", c)
}

// SetSelfInfo
// @Tags      SysUser
// @Summary   设置自身信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, 用户名, 昵称, 头像链接"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "设置用户信息"
// @Router    /api/v1/user/self/info/set [post]
func (b *BaseApi) SetSelfInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user.ID = utils.GetUserID(c)
	err = userService.SetSelfInfo(system.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		Enable:    user.Enable,
	})
	if err != nil {
		utils.HandleError(c, err, "设置失败")
		return
	}
	response.OkWithMessage("设置成功", c)
}

// SetSelfSetting
// @Tags      SysUser
// @Summary   设置用户配置
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      map[string]interface{}  true  "用户配置数据"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "设置用户配置"
// @Router    /api/v1/user/self/setting/set [post]
func (b *BaseApi) SetSelfSetting(c *gin.Context) {
	var req common.JSONMap
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = userService.SetSelfSetting(req, utils.GetUserID(c))
	if err != nil {
		utils.HandleError(c, err, "设置失败")
		return
	}
	response.OkWithMessage("设置成功", c)
}

// GetUserInfo
// @Tags      SysUser
// @Summary   获取当前登录用户信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "获取用户信息"
// @Router    /api/v1/user/info [get]
func (a *BaseApi) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	userResp, err := userService.GetUserInfo(uuid)
	if err != nil {
		utils.HandleError(c, err, "获取失败")
		return
	}
	response.OkWithDetailed(gin.H{"userInfo": userResp}, "获取成功", c)
}

// ResetPassword
// @Tags      SysUser
// @Summary   重置用户密码
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      systemReq.ResetPassword        true  "ID, Password"
// @Success   200   {object}  response.Response{msg=string}  "重置用户密码"
// @Router    /api/v1/user/password/reset [post]
func (b *BaseApi) ResetPassword(c *gin.Context) {
	var rps systemReq.ResetPassword
	err := c.ShouldBindJSON(&rps)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.ResetPassword(rps.ID, rps.Password)
	if err != nil {
		utils.HandleError(c, err, "重置失败")
		return
	}
	response.OkWithMessage("重置成功", c)
}
