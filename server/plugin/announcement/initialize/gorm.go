package initialize

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/plugin/announcement/model"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

func Gorm(ctx context.Context) {
	err := global.GVA_DB.WithContext(ctx).AutoMigrate(
		new(model.Info),
	)
	if err != nil {
		err = errors.Wrap(err, "注册表失败!")
		logx.Error(fmt.Sprintf("%+v", err))
	}
}
