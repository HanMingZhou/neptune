package system

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model/common"
	orderModel "gin-vue-admin/model/order"
	"gin-vue-admin/model/system"
	systemReq "gin-vue-admin/model/system/request"
	systemRes "gin-vue-admin/model/system/response"
	"gin-vue-admin/utils"
	helper "gin-vue-admin/utils/k8s"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserService struct{}

var UserServiceApp = new(UserService)

func (userService *UserService) Register(req systemReq.Register) (userResp systemRes.UserResponse, err error) {
	var user system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", req.Username).First(&user).Error, gorm.ErrRecordNotFound) {
		return userResp, errors.New("用户名已注册")
	}

	nickName := req.NickName
	if nickName == "" {
		nickName = req.Username
	}

	u := system.SysUser{
		Username:    req.Username,
		NickName:    nickName,
		Password:    utils.BcryptHash(req.Password),
		HeaderImg:   req.HeaderImg,
		AuthorityId: system.DefaultAuthorityId,
		Authorities: []system.SysAuthority{{AuthorityId: system.DefaultAuthorityId}},
		Enable:      1,
		Phone:       req.Phone,
		Email:       req.Email,
		UUID:        uuid.New(),
		Namespace:   helper.GenerateNamespaceName(uuid.New().String()),
	}
	err = global.GVA_DB.Create(&u).Error
	if err != nil {
		return userResp, err
	}

	// 为新用户创建钱包
	wallet := orderModel.Wallet{UserId: u.ID, Balance: 0}
	global.GVA_DB.FirstOrCreate(&wallet, orderModel.Wallet{UserId: u.ID})

	return userService.MapToResponse(u), nil
}

func (userService *UserService) Login(username, password string) (userResp systemRes.UserResponse, err error) {
	if nil == global.GVA_DB {
		return userResp, errors.Errorf("db not init")
	}

	var user system.SysUser
	err = global.GVA_DB.Where("username = ?", username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err != nil {
		return userResp, errors.New("用户名不存在或密码错误")
	}

	if ok := utils.BcryptCheck(password, user.Password); !ok {
		return userResp, errors.New("密码错误")
	}

	MenuServiceApp.UserAuthorityDefaultRouter(&user)
	return userService.MapToResponse(user), nil
}

func (userService *UserService) ChangePassword(id uint, oldPassword, newPassword string) (err error) {
	var user system.SysUser
	err = global.GVA_DB.Select("id, password").Where("id = ?", id).First(&user).Error
	if err != nil {
		return errors.New("用户不存在")
	}
	if ok := utils.BcryptCheck(oldPassword, user.Password); !ok {
		return errors.New("原密码错误")
	}
	pwd := utils.BcryptHash(newPassword)
	return global.GVA_DB.Model(&user).Update("password", pwd).Error
}

func (userService *UserService) GetUserInfoList(info systemReq.GetUserList) (list []systemRes.UserResponse, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysUser{})
	var userList []system.SysUser

	if info.NickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+info.NickName+"%")
	}
	if info.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+info.Phone+"%")
	}
	if info.Username != "" {
		db = db.Where("username LIKE ?", "%"+info.Username+"%")
	}
	if info.Email != "" {
		db = db.Where("email LIKE ?", "%"+info.Email+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return nil, total, err
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").Find(&userList).Error
	if err != nil {
		return nil, total, err
	}

	resps := make([]systemRes.UserResponse, 0, len(userList))
	for _, u := range userList {
		resps = append(resps, userService.MapToResponse(u))
	}

	return resps, total, nil
}

func (userService *UserService) SetUserAuthority(id uint, authorityId uint) (err error) {
	assignErr := global.GVA_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&system.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}

	var authority system.SysAuthority
	err = global.GVA_DB.Where("authority_id = ?", authorityId).First(&authority).Error
	if err != nil {
		return errors.New("角色不存在")
	}

	err = global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", id).Update("authority_id", authorityId).Error
	return err
}

func (userService *UserService) SetUserAuthorities(adminAuthorityID, id uint, authorityIds []uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var user system.SysUser
		if err := tx.Where("id = ?", id).First(&user).Error; err != nil {
			return errors.New("查询用户数据失败")
		}
		if err := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error; err != nil {
			return err
		}
		var useAuthority []system.SysUserAuthority
		for _, v := range authorityIds {
			if err := AuthorityServiceApp.CheckAuthorityIDAuth(adminAuthorityID, v); err != nil {
				return err
			}
			useAuthority = append(useAuthority, system.SysUserAuthority{
				SysUserId: id, SysAuthorityAuthorityId: v,
			})
		}
		if err := tx.Create(&useAuthority).Error; err != nil {
			return err
		}
		return tx.Model(&user).Update("authority_id", authorityIds[0]).Error
	})
}

func (userService *UserService) DeleteUser(id int) (err error) {
	var user system.SysUser
	if err = global.GVA_DB.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&system.SysUser{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error; err != nil {
			return err
		}
		if global.GVA_CONFIG.System.UseRedis {
			_ = JwtServiceApp.ClearRedisJWTs(user.Username)
		}
		return nil
	})
}

func (userService *UserService) SetUserInfo(req system.SysUser) error {
	err := global.GVA_DB.Model(&system.SysUser{}).
		Select("updated_at", "nick_name", "header_img", "phone", "email", "enable").
		Where("id=?", req.ID).
		Updates(map[string]interface{}{
			"updated_at": time.Now(),
			"nick_name":  req.NickName,
			"header_img": req.HeaderImg,
			"phone":      req.Phone,
			"email":      req.Email,
			"enable":     req.Enable,
		}).Error
	if err != nil {
		return err
	}

	if req.Enable == 2 && global.GVA_CONFIG.System.UseRedis {
		var user system.SysUser
		if err := global.GVA_DB.First(&user, req.ID).Error; err == nil {
			_ = JwtServiceApp.ClearRedisJWTs(user.Username)
		}
	}
	return nil
}

func (userService *UserService) SetSelfInfo(req system.SysUser) error {
	return global.GVA_DB.Model(&system.SysUser{}).Where("id=?", req.ID).Updates(req).Error
}

func (userService *UserService) SetSelfSetting(req common.JSONMap, uid uint) error {
	return global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", uid).Update("origin_setting", req).Error
}

func (userService *UserService) GetUserInfo(uuid uuid.UUID) (userResp systemRes.UserResponse, err error) {
	var user system.SysUser
	err = global.GVA_DB.Preload("Authorities").Preload("Authority").First(&user, "uuid = ?", uuid).Error
	if err != nil {
		return userResp, err
	}
	MenuServiceApp.UserAuthorityDefaultRouter(&user)
	return userService.MapToResponse(user), nil
}

func (userService *UserService) ResetPassword(ID uint, password string) (err error) {
	return global.GVA_DB.Model(&system.SysUser{}).Where("id = ?", ID).Update("password", utils.BcryptHash(password)).Error
}

func (userService *UserService) MapToResponse(user system.SysUser) systemRes.UserResponse {
	return systemRes.UserResponse{
		ID:          user.ID,
		UUID:        user.UUID.String(),
		Username:    user.Username,
		NickName:    user.NickName,
		HeaderImg:   user.HeaderImg,
		AuthorityId: user.AuthorityId,
		Authority:   user.Authority,
		Authorities: user.Authorities,
		Phone:       user.Phone,
		Email:       user.Email,
		Enable:      user.Enable,
		Namespace:   user.Namespace,
	}
}

func (userService *UserService) FindUserByUuid(uuid string) (user system.SysUser, err error) {
	if err = global.GVA_DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
