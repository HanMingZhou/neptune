package initialize

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/plugin/announcement/plugin"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

func Viper() {
	err := global.GVA_VP.UnmarshalKey("announcement", &plugin.Config)
	if err != nil {
		err = errors.Wrap(err, "初始化配置文件失败!")
		logx.Error(fmt.Sprintf("%+v", err))
	}
}
