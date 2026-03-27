package system

import (
	"context"
	"gin-vue-admin/global"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/service/system"

	"github.com/pkg/errors"
)

type initApiIgnore struct{}

const initOrderApiIgnore = initOrderApi + 1

// auto run
func init() {
	system.RegisterInit(initOrderApiIgnore, &initApiIgnore{})
}

func (i *initApiIgnore) InitializerName() string {
	return sysModel.SysIgnoreApi{}.TableName()
}

func (i *initApiIgnore) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, global.GVA_DB.AutoMigrate(&sysModel.SysIgnoreApi{})
}

func (i *initApiIgnore) TableCreated(ctx context.Context) bool {
	return global.GVA_DB.Migrator().HasTable(&sysModel.SysIgnoreApi{})
}

func (i *initApiIgnore) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []sysModel.SysIgnoreApi{
		{Method: "GET", Path: "/swagger/*any"},
		{Method: "GET", Path: "/api/v1/freshCasbin"},
		{Method: "GET", Path: "/uploads/file/*filepath"},
		{Method: "GET", Path: "/health"},
		{Method: "HEAD", Path: "/uploads/file/*filepath"},
		{Method: "POST", Path: "/api/v1/autocode/llm/add"},
		{Method: "POST", Path: "/api/v1/system/reload"},
		{Method: "POST", Path: "/api/v1/base/login"},
		{Method: "POST", Path: "/api/v1/base/captcha"},
		{Method: "POST", Path: "/api/v1/init/add"},
		{Method: "POST", Path: "/api/v1/init/check"},
		{Method: "GET", Path: "/api/v1/info/get/datasource"},
		{Method: "GET", Path: "/api/v1/info/get/public"},
	}
	// 批量查询已存在的 API Ignore，避免循环执行 SQL
	var existingApis []sysModel.SysIgnoreApi
	if err := global.GVA_DB.Select("path", "method").Find(&existingApis).Error; err != nil {
		return ctx, errors.Wrap(err, "查询已存在API Ignore失败")
	}

	existingMap := make(map[string]bool)
	for _, a := range existingApis {
		existingMap[a.Path+":"+a.Method] = true
	}

	var missingEntities []sysModel.SysIgnoreApi
	for _, entity := range entities {
		if !existingMap[entity.Path+":"+entity.Method] {
			missingEntities = append(missingEntities, entity)
		}
	}

	if len(missingEntities) > 0 {
		if err := global.GVA_DB.Create(&missingEntities).Error; err != nil {
			return ctx, errors.Wrap(err, sysModel.SysIgnoreApi{}.TableName()+"表数据批量初始化失败!")
		}
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
