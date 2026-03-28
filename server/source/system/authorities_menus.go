package system

import (
	"context"
	"gin-vue-admin/global"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/service/system"

	"github.com/pkg/errors"
)

const initOrderMenuAuthority = initOrderMenu + initOrderAuthority

type initMenuAuthority struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}

func (i *initMenuAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initMenuAuthority) TableCreated(ctx context.Context) bool {
	return false // always replace
}

func (i *initMenuAuthority) InitializerName() string {
	return "sys_menu_authorities"
}

func (i *initMenuAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	authorities := []sysModel.SysAuthority{}
	if err := global.GVA_DB.Order("authority_id").Find(&authorities).Error; err != nil {
		return ctx, errors.Wrap(system.ErrMissingDependentContext, "创建 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}

	// 始终从数据库重新加载所有菜单，确保获取最新数据
	var allMenus []sysModel.SysBaseMenu
	if err := global.GVA_DB.Order("sort").Find(&allMenus).Error; err != nil {
		return ctx, errors.Wrap(err, "创建 [菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}
	next = ctx

	// 构建角色映射
	authMap := make(map[uint]sysModel.SysAuthority)
	for _, auth := range authorities {
		authMap[auth.AuthorityId] = auth
	}

	adminAuth := authMap[888]
	testAuth := authMap[9528]
	userAuth := authMap[8881]

	// 1. 超级管理员角色(888) - 拥有所有菜单权限
	if err = global.GVA_DB.Model(&adminAuth).Association("SysBaseMenus").Replace(allMenus); err != nil {
		return next, errors.Wrap(err, "为超级管理员分配菜单失败")
	}

	// 2. 测试角色(9528) - 拥有所有菜单权限
	if err = global.GVA_DB.Model(&testAuth).Association("SysBaseMenus").Replace(allMenus); err != nil {
		return next, errors.Wrap(err, "为测试角色分配菜单失败")
	}

	// 3. 普通用户角色(8881) - 拥有受限菜单权限
	var menu8881 []sysModel.SysBaseMenu

	// 定义普通用户可访问的顶级菜单名称
	userTopMenuNames := map[string]bool{
		"dashboard": true, // 控制台概览
		"notebooks": true, // 容器实例
		"training":  true, // 训练任务
		"inference": true, // 推理服务
		"sshkeys":   true, // SSH 密钥
		"storage":   true, // 存储卷
		"images":    true, // 镜像管理
		"order":     true, // 费用账单（父菜单）
		"account":   true, // 个人中心（父菜单）
	}

	// 定义普通用户可访问的子菜单名称
	userSubMenuNames := map[string]bool{
		// account 子菜单
		"security":      true, // 账号安全
		"accessRecords": true, // 访问记录
		"person":        true, // 个人信息
		// order 子菜单
		"transactions": true, // 交易记录
		"usage":        true, // 资源用量
		"invoice":      true, // 发票管理
	}

	for _, menu := range allMenus {
		if menu.ParentId == 0 {
			// 顶级菜单
			if userTopMenuNames[menu.Name] {
				menu8881 = append(menu8881, menu)
			}
		} else {
			// 子菜单：检查是否在允许列表中
			if userSubMenuNames[menu.Name] {
				menu8881 = append(menu8881, menu)
			}
		}
	}

	if err = global.GVA_DB.Model(&userAuth).Association("SysBaseMenus").Replace(menu8881); err != nil {
		return next, errors.Wrap(err, "为普通用户分配菜单失败")
	}

	var transactionsMenu sysModel.SysBaseMenu
	if err = global.GVA_DB.Where("name = ?", "transactions").First(&transactionsMenu).Error; err != nil {
		return next, errors.Wrap(err, "查询交易记录菜单失败")
	}

	var transactionBtns []sysModel.SysBaseMenuBtn
	if err = global.GVA_DB.Where("sys_base_menu_id = ?", transactionsMenu.ID).Find(&transactionBtns).Error; err != nil {
		return next, errors.Wrap(err, "查询交易记录菜单按钮失败")
	}

	btnIDByName := make(map[string]uint, len(transactionBtns))
	for _, btn := range transactionBtns {
		btnIDByName[btn.Name] = btn.ID
	}

	assignBtn := func(authorityId uint, btnName string) error {
		btnID, ok := btnIDByName[btnName]
		if !ok {
			return nil
		}
		return global.GVA_DB.FirstOrCreate(&sysModel.SysAuthorityBtn{}, sysModel.SysAuthorityBtn{
			AuthorityId:      authorityId,
			SysMenuID:        transactionsMenu.ID,
			SysBaseMenuBtnID: btnID,
		}).Error
	}

	defaultBtnAssignments := map[uint][]string{
		adminAuth.AuthorityId: {"recharge_system"},
		testAuth.AuthorityId:  {"recharge_alipay", "recharge_wechat"},
		userAuth.AuthorityId:  {"recharge_alipay", "recharge_wechat"},
	}
	for authorityId, btnNames := range defaultBtnAssignments {
		for _, btnName := range btnNames {
			if err = assignBtn(authorityId, btnName); err != nil {
				return next, errors.Wrap(err, "初始化交易记录按钮权限失败")
			}
		}
	}

	return next, nil
}
