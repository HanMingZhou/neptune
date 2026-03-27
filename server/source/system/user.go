package system

import (
	"context"
	"gin-vue-admin/global"
	orderModel "gin-vue-admin/model/order"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/service/system"
	"gin-vue-admin/utils"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderUser = initOrderAuthority + 1

type initUser struct{}

// auto run
func init() {
	system.RegisterInit(initOrderUser, &initUser{})
}

func (i *initUser) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, global.GVA_DB.AutoMigrate(&sysModel.SysUser{})
}

func (i *initUser) TableCreated(ctx context.Context) bool {
	return global.GVA_DB.Migrator().HasTable(&sysModel.SysUser{})
}

func (i *initUser) InitializerName() string {
	return sysModel.SysUser{}.TableName()
}

func (i *initUser) InitializeData(ctx context.Context) (next context.Context, err error) {
	entities := []sysModel.SysUser{
		{
			UUID:        uuid.New(),
			Username:    "admin",
			Password:    utils.BcryptHash("123456"),
			NickName:    "Mr.奇淼",
			HeaderImg:   "https://qmplusimg.henrongyi.top/gva_header.jpg",
			AuthorityId: sysModel.SysAuthorityId888,
			Phone:       "17611111111",
			Email:       "333333333@qq.com",
		},
		{
			UUID:        uuid.New(),
			Username:    "a303176530",
			Password:    utils.BcryptHash("123456"),
			NickName:    "用户1",
			HeaderImg:   "https://qmplusimg.henrongyi.top/1572075907logo.png",
			AuthorityId: sysModel.SysAuthorityId9528,
			Phone:       "17611111111",
			Email:       "333333333@qq.com"},
	}
	// 批量查询已存在的 User，并把数据回填到 entities 中以获取 ID
	for k, entity := range entities {
		var existing sysModel.SysUser
		if !errors.Is(global.GVA_DB.Where("username = ?", entity.Username).First(&existing).Error, gorm.ErrRecordNotFound) {
			entities[k] = existing
		}
	}

	var missingEntities []sysModel.SysUser
	for _, entity := range entities {
		if entity.ID == 0 {
			missingEntities = append(missingEntities, entity)
		}
	}

	if len(missingEntities) > 0 {
		if err = global.GVA_DB.Create(&missingEntities).Error; err != nil {
			return ctx, errors.Wrap(err, sysModel.SysUser{}.TableName()+"表数据批量初始化失败!")
		}
		// 重新写回到 entities 以确保关联操作有 ID
		for k, entity := range entities {
			if entity.ID == 0 {
				global.GVA_DB.Where("username = ?", entity.Username).First(&entities[k])
			}
		}
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	authorityEntities, ok := ctx.Value(new(initAuthority).InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return next, errors.Wrap(system.ErrMissingDependentContext, "创建 [用户-权限] 关联失败, 未找到权限表初始化数据")
	}
	if err = global.GVA_DB.Model(&entities[0]).Association("Authorities").Replace(authorityEntities); err != nil {
		return next, err
	}
	if err = global.GVA_DB.Model(&entities[1]).Association("Authorities").Replace(authorityEntities[:1]); err != nil {
		return next, err
	}

	// 为每个初始化用户创建钱包并充值 1000000.0
	for _, user := range entities {
		if user.ID == 0 {
			continue
		}
		wallet := orderModel.Wallet{UserId: user.ID}
		result := global.GVA_DB.Where("user_id = ?", user.ID).FirstOrCreate(&wallet)
		if result.Error != nil {
			return next, errors.Wrap(result.Error, "初始化用户钱包失败")
		}
		if result.RowsAffected > 0 {
			// 新创建的钱包，设置初始余额
			global.GVA_DB.Model(&wallet).Update("balance", 1000000.0)
		}
	}

	return next, err
}
