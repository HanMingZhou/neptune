package system

import (
	"context"
	"gin-vue-admin/global"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/service/system"

	"github.com/pkg/errors"
)

type initExcelTemplate struct{}

const initOrderExcelTemplate = initOrderDictDetail + 1

// auto run
func init() {
	system.RegisterInit(initOrderExcelTemplate, &initExcelTemplate{})
}

func (i *initExcelTemplate) InitializerName() string {
	return "sys_export_templates"
}

func (i *initExcelTemplate) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, global.GVA_DB.AutoMigrate(&sysModel.SysExportTemplate{})
}

func (i *initExcelTemplate) TableCreated(ctx context.Context) bool {
	return global.GVA_DB.Migrator().HasTable(&sysModel.SysExportTemplate{})
}

func (i *initExcelTemplate) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []sysModel.SysExportTemplate{
		{
			Name:       "api",
			TableName:  "sys_apis",
			TemplateID: "api",
			TemplateInfo: `{
"path":"路径",
"method":"方法（大写）",
"description":"方法介绍",
"api_group":"方法分组"
}`,
		},
	}
	// 批量查询已存在的 Excel Template，避免循环执行 SQL
	var existingTemplates []sysModel.SysExportTemplate
	if err := global.GVA_DB.Select("template_id").Find(&existingTemplates).Error; err != nil {
		return ctx, errors.Wrap(err, "查询已存在Excel Template失败")
	}

	existingMap := make(map[string]bool)
	for _, t := range existingTemplates {
		existingMap[t.TemplateID] = true
	}

	var missingEntities []sysModel.SysExportTemplate
	for _, entity := range entities {
		if !existingMap[entity.TemplateID] {
			missingEntities = append(missingEntities, entity)
		}
	}

	if len(missingEntities) > 0 {
		if err := global.GVA_DB.Create(&missingEntities).Error; err != nil {
			return ctx, errors.Wrap(err, "sys_export_templates"+"表数据批量初始化失败!")
		}
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
