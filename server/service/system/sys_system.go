package system

import (
	"gin-vue-admin/config"
	"gin-vue-admin/global"
	"gin-vue-admin/model/system"
	"gin-vue-admin/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSystemConfig
//@description: 读取配置文件
//@return: conf config.Server, err error

type SystemConfigService struct{}

var SystemConfigServiceApp = new(SystemConfigService)

func (systemConfigService *SystemConfigService) GetSystemConfig() (conf config.Server, err error) {
	return global.GVA_CONFIG, nil
}

// @description   set system config,
//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetSystemConfig
//@description: 设置配置文件
//@param: system model.System
//@return: err error

func (systemConfigService *SystemConfigService) SetSystemConfig(system system.System) (err error) {
	cs := utils.StructToMap(system.Config)
	for k, v := range cs {
		global.GVA_VP.Set(k, v)
	}
	err = global.GVA_VP.WriteConfig()
	return err
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: GetServerInfo
//@description: 获取服务器信息
//@return: server *utils.Server, err error

func (systemConfigService *SystemConfigService) GetServerInfo() (server *utils.Server, err error) {
	var s utils.Server
	s.Os = utils.InitOS()
	if s.Cpu, err = utils.InitCPU(); err != nil {
		logx.Error("func utils.InitCPU() Failed", err)
		return &s, err
	}
	if s.Ram, err = utils.InitRAM(); err != nil {
		logx.Error("func utils.InitRAM() Failed", err)
		return &s, err
	}
	if s.Disk, err = utils.InitDisk(); err != nil {
		logx.Error("func utils.InitDisk() Failed", err)
		return &s, err
	}

	return &s, nil
}
