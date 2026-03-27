package initialize

import (
	"gin-vue-admin/global"
	"gin-vue-admin/service/system"
	// 导入 source/system 包以注册初始化器
	_ "gin-vue-admin/source/system"

	"github.com/zeromicro/go-zero/core/logx"
)

// InitBaseData 在项目启动时检测并插入基础数据
// 直接复用 source/system 包中已注册的初始化器
func InitBaseData() {
	if global.GVA_DB == nil {
		logx.Error("数据库未连接，跳过基础数据初始化")
		return
	}

	initDBService := system.InitDBService{}
	if err := initDBService.InitDataOnly(); err != nil {
		logx.Error("基础数据初始化失败", err)
		return
	}

	logx.Info("基础数据初始化完成！")
	logx.Info("默认管理员账号: admin, 密码: 123456")
}
