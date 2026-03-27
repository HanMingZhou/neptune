package utils

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model/system"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

func RegisterApis(apis ...system.SysApi) {
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, api := range apis {
			err := tx.Model(system.SysApi{}).Where("path = ? AND method = ? AND api_group = ? ", api.Path, api.Method, api.ApiGroup).FirstOrCreate(&api).Error
			if err != nil {
				logx.Error("注册API失败", logx.Field("err", err), logx.Field("api", api.Path), logx.Field("method", api.Method), logx.Field("apiGroup", api.ApiGroup))
				return err
			}
		}
		return nil
	})
	if err != nil {
		logx.Error("注册API失败", logx.Field("err", err))
	}
}

func RegisterMenus(menus ...system.SysBaseMenu) {
	parentMenu := menus[0]
	otherMenus := menus[1:]
	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(system.SysBaseMenu{}).Where("name = ? ", parentMenu.Name).FirstOrCreate(&parentMenu).Error
		if err != nil {
			logx.Error("注册菜单失败", logx.Field("err", err), logx.Field("menu", parentMenu.Name))
			return errors.Wrap(err, "注册菜单失败")
		}
		pid := parentMenu.ID
		for i := range otherMenus {
			otherMenus[i].ParentId = pid
			err = tx.Model(system.SysBaseMenu{}).Where("name = ? ", otherMenus[i].Name).FirstOrCreate(&otherMenus[i]).Error
			if err != nil {
				logx.Error("注册菜单失败", logx.Field("err", err), logx.Field("menu", otherMenus[i].Name))
				return errors.Wrap(err, "注册菜单失败")
			}
		}

		return nil
	})
	if err != nil {
		logx.Error("注册菜单失败", logx.Field("err", err))
	}

}
