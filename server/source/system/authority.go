package system

import (
	"context"
	"gin-vue-admin/global"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/service/system"
	"gin-vue-admin/utils"

	"github.com/pkg/errors"
)

const initOrderAuthority = initOrderCasbin + 1

type initAuthority struct{}

// auto run
func init() {
	system.RegisterInit(initOrderAuthority, &initAuthority{})
}

func (i *initAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, global.GVA_DB.AutoMigrate(&sysModel.SysAuthority{})
}

func (i *initAuthority) TableCreated(ctx context.Context) bool {
	return global.GVA_DB.Migrator().HasTable(&sysModel.SysAuthority{})
}

func (i *initAuthority) InitializerName() string {
	return sysModel.SysAuthority{}.TableName()
}

func (i *initAuthority) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []sysModel.SysAuthority{
		{AuthorityId: sysModel.SysAuthorityId888, AuthorityName: "管理员", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: sysModel.SysAuthorityId9528, AuthorityName: "测试角色", ParentId: utils.Pointer[uint](0), DefaultRouter: "dashboard"},
		{AuthorityId: 8881, AuthorityName: "普通用户", ParentId: utils.Pointer[uint](sysModel.SysAuthorityId888), DefaultRouter: "dashboard"},
	}

	// 批量查询已存在的 Authority，避免循环执行 SQL
	var existingAuthorities []sysModel.SysAuthority
	if err := global.GVA_DB.Select("authority_id").Find(&existingAuthorities).Error; err != nil {
		return ctx, errors.Wrap(err, "查询已存在Authority失败")
	}

	existingMap := make(map[uint]bool)
	for _, a := range existingAuthorities {
		existingMap[a.AuthorityId] = true
	}

	var missingEntities []sysModel.SysAuthority
	for _, entity := range entities {
		if !existingMap[entity.AuthorityId] {
			missingEntities = append(missingEntities, entity)
		}
	}

	if len(missingEntities) > 0 {
		if err := global.GVA_DB.Create(&missingEntities).Error; err != nil {
			return ctx, errors.Wrapf(err, "%s表数据初始化失败!", sysModel.SysAuthority{}.TableName())
		}
	}
	// data authority
	if err := global.GVA_DB.Model(&entities[0]).Association("DataAuthorityId").Replace(
		[]*sysModel.SysAuthority{
			{AuthorityId: sysModel.SysAuthorityId888},
			{AuthorityId: sysModel.SysAuthorityId9528},
			{AuthorityId: 8881},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			global.GVA_DB.Model(&entities[0]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}
	if err := global.GVA_DB.Model(&entities[1]).Association("DataAuthorityId").Replace(
		[]*sysModel.SysAuthority{
			{AuthorityId: sysModel.SysAuthorityId9528},
			{AuthorityId: 8881},
		}); err != nil {
		return ctx, errors.Wrapf(err, "%s表数据初始化失败!",
			global.GVA_DB.Model(&entities[1]).Association("DataAuthorityId").Relationship.JoinTable.Name)
	}

	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
