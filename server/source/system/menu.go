package system

import (
	"context"
	"gin-vue-admin/global"
	"gin-vue-admin/service/system"

	. "gin-vue-admin/model/system"
	"github.com/pkg/errors"
)

const initOrderMenu = initOrderAuthority + 1

type initMenu struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenu, &initMenu{})
}

func (i *initMenu) InitializerName() string {
	return SysBaseMenu{}.TableName()
}

func (i *initMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, global.GVA_DB.AutoMigrate(
		&SysBaseMenu{},
		&SysBaseMenuParameter{},
		&SysBaseMenuBtn{},
	)
}

func (i *initMenu) TableCreated(ctx context.Context) bool {
	m := global.GVA_DB.Migrator()
	return m.HasTable(&SysBaseMenu{}) &&
		m.HasTable(&SysBaseMenuParameter{}) &&
		m.HasTable(&SysBaseMenuBtn{})
}

func (i *initMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	// ============================================================
	// 菜单结构 - 基于前端 constants.ts (NAVIGATION_GROUPS)
	// ============================================================

	// 1. 顶级菜单
	level0 := []SysBaseMenu{
		// compute 分组
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "dashboard", Name: "dashboard", Component: "view/console/dashboard/index.vue", Sort: 1, Meta: Meta{Title: "控制台概览", Icon: "dashboard"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "notebooks", Name: "notebooks", Component: "view/console/notebooks/index.vue", Sort: 2, Meta: Meta{Title: "容器实例", Icon: "terminal"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "training", Name: "training", Component: "view/console/training/index.vue", Sort: 3, Meta: Meta{Title: "训练任务", Icon: "model_training"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "inference", Name: "inference", Component: "view/console/inference/index.vue", Sort: 4, Meta: Meta{Title: "推理服务", Icon: "rocket_launch"}},

		// resources 分组
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "sshkeys", Name: "sshkeys", Component: "view/console/sshkeys/index.vue", Sort: 5, Meta: Meta{Title: "SSH 密钥", Icon: "vpn_key"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "storage", Name: "storage", Component: "view/console/storage/index.vue", Sort: 6, Meta: Meta{Title: "存储卷", Icon: "storage"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "images", Name: "images", Component: "view/console/images/index.vue", Sort: 7, Meta: Meta{Title: "镜像管理", Icon: "photo_library"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "order", Name: "order", Component: "view/routerHolder.vue", Sort: 8, Meta: Meta{Title: "费用账单", Icon: "payments"}},

		// management 分组 - account 是父菜单
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "account", Name: "account", Component: "view/routerHolder.vue", Sort: 9, Meta: Meta{Title: "个人中心", Icon: "person"}},

		// admin 分组 - 管理员菜单
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "admin", Name: "admin", Component: "view/routerHolder.vue", Sort: 100, Meta: Meta{Title: "系统管理", Icon: "admin_panel_settings"}},
	}

	var existingMenus []SysBaseMenu
	if err := global.GVA_DB.Select("id", "name").Find(&existingMenus).Error; err != nil {
		return ctx, errors.Wrap(err, "查询已存在菜单失败")
	}
	existingMap := make(map[string]uint)
	for _, m := range existingMenus {
		existingMap[m.Name] = m.ID
	}

	createMissing := func(menus []SysBaseMenu) error {
		var missing []SysBaseMenu
		for k, m := range menus {
			if id, ok := existingMap[m.Name]; ok {
				menus[k].ID = id
			} else {
				missing = append(missing, menus[k])
			}
		}
		if len(missing) > 0 {
			if err := global.GVA_DB.Create(&missing).Error; err != nil {
				return err
			}
			// 刷新 ID
			for _, m := range missing {
				existingMap[m.Name] = m.ID
			}
		}
		return nil
	}

	if err = createMissing(level0); err != nil {
		return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"顶级菜单初始化失败!")
	}

	// 获取父菜单 ID
	var accountMenu, adminMenu, orderMenu SysBaseMenu
	if err = global.GVA_DB.Where("name = ?", "account").First(&accountMenu).Error; err != nil {
		return ctx, errors.Wrap(err, "查询 account 菜单失败")
	}
	if err = global.GVA_DB.Where("name = ?", "admin").First(&adminMenu).Error; err != nil {
		return ctx, errors.Wrap(err, "查询 admin 菜单失败")
	}
	if err = global.GVA_DB.Where("name = ?", "order").First(&orderMenu).Error; err != nil {
		return ctx, errors.Wrap(err, "查询 order 菜单失败")
	}

	// 2. account 子菜单
	accountChildren := []SysBaseMenu{
		{MenuLevel: 1, Hidden: false, ParentId: accountMenu.ID, Path: "security", Name: "security", Component: "view/console/account/security.vue", Sort: 1, Meta: Meta{Title: "账号安全", Icon: "shield"}},
		{MenuLevel: 1, Hidden: false, ParentId: accountMenu.ID, Path: "records", Name: "accessRecords", Component: "view/console/account/accesslog.vue", Sort: 2, Meta: Meta{Title: "访问记录", Icon: "history"}},
		{MenuLevel: 1, Hidden: false, ParentId: accountMenu.ID, Path: "person", Name: "person", Component: "view/console/account/person.vue", Sort: 3, Meta: Meta{Title: "个人信息", Icon: "user-filled"}},
	}

	if err = createMissing(accountChildren); err != nil {
		return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"account子菜单初始化失败!")
	}
	// 2.5 order 子菜单
	orderChildren := []SysBaseMenu{
		{MenuLevel: 1, Hidden: false, ParentId: orderMenu.ID, Path: "transactions", Name: "transactions", Component: "view/console/order/transactions.vue", Sort: 1, Meta: Meta{Title: "交易记录", Icon: "receipt"}},
		{MenuLevel: 1, Hidden: false, ParentId: orderMenu.ID, Path: "usage", Name: "usage", Component: "view/console/order/usage.vue", Sort: 2, Meta: Meta{Title: "资源用量", Icon: "data_usage"}},
		{MenuLevel: 1, Hidden: false, ParentId: orderMenu.ID, Path: "invoice", Name: "invoice", Component: "view/console/order/invoice.vue", Sort: 3, Meta: Meta{Title: "发票管理", Icon: "description"}},
	}
	if err = createMissing(orderChildren); err != nil {
		return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"order子菜单初始化失败!")
	}

	// 3. admin 子菜单 - 扁平化结构（避免3层嵌套）
	adminChildren := []SysBaseMenu{
		// 产品管理
		{MenuLevel: 1, Hidden: false, ParentId: adminMenu.ID, Path: "products", Name: "products", Component: "view/superAdmin/cms/product.vue", Sort: 1, Meta: Meta{Title: "产品管理", Icon: "inventory_2"}},
		// 角色管理
		{MenuLevel: 1, Hidden: false, ParentId: adminMenu.ID, Path: "roles", Name: "roles", Component: "view/superAdmin/authority/authority.vue", Sort: 2, Meta: Meta{Title: "角色管理", Icon: "badge"}},
		// 菜单管理
		{MenuLevel: 1, Hidden: false, ParentId: adminMenu.ID, Path: "menus", Name: "menus", Component: "view/superAdmin/menu/menu.vue", Sort: 3, Meta: Meta{Title: "菜单管理", Icon: "menu_open", KeepAlive: true}},
		// 接口管理
		{MenuLevel: 1, Hidden: false, ParentId: adminMenu.ID, Path: "apis", Name: "apis", Component: "view/superAdmin/api/api.vue", Sort: 4, Meta: Meta{Title: "接口管理", Icon: "api", KeepAlive: true}},
		// 用户管理
		{MenuLevel: 1, Hidden: false, ParentId: adminMenu.ID, Path: "users", Name: "users", Component: "view/superAdmin/user/user.vue", Sort: 5, Meta: Meta{Title: "用户管理", Icon: "group"}},
		// 操作审计
		{MenuLevel: 1, Hidden: false, ParentId: adminMenu.ID, Path: "operations", Name: "operations", Component: "view/superAdmin/operation/sysOperationRecord.vue", Sort: 6, Meta: Meta{Title: "操作审计", Icon: "receipt_long"}},
		// 集群管理
		{MenuLevel: 1, Hidden: false, ParentId: adminMenu.ID, Path: "clusters", Name: "clusterManage", Component: "view/superAdmin/cms/cluster.vue", Sort: 7, Meta: Meta{Title: "集群管理", Icon: "cloud"}},
		// 节点管理
		{MenuLevel: 1, Hidden: false, ParentId: adminMenu.ID, Path: "nodes", Name: "nodeManage", Component: "view/superAdmin/cms/node.vue", Sort: 8, Meta: Meta{Title: "节点管理", Icon: "hub"}},
	}

	if err = createMissing(adminChildren); err != nil {
		return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"admin子菜单初始化失败!")
	}

	// 组合所有菜单作为返回结果
	allEntities := append(level0, accountChildren...)
	allEntities = append(allEntities, orderChildren...)
	allEntities = append(allEntities, adminChildren...)
	// rolesChildren 已移除（扁平化到 adminChildren 中）

	next = context.WithValue(ctx, i.InitializerName(), allEntities)
	return next, nil
}
