package system

import (
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

type MenuRouter struct{}

func (s *MenuRouter) InitMenuRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	menuRouter := Router.Group("menu").Use(middleware.OperationRecord())
	{
		menuRouter.POST("add", authorityMenuApi.AddBaseMenu)                   // 新增菜单
		menuRouter.POST("authority/update", authorityMenuApi.AddMenuAuthority) //	增加menu和角色关联关系
		menuRouter.POST("delete", authorityMenuApi.DeleteBaseMenu)             // 删除菜单
		menuRouter.POST("update", authorityMenuApi.UpdateBaseMenu)             // 更新菜单
		menuRouter.POST("get", authorityMenuApi.GetMenu)                       // 获取菜单树
		menuRouter.POST("list", authorityMenuApi.GetMenuList)                  // 分页获取基础menu列表
		menuRouter.POST("tree", authorityMenuApi.GetBaseMenuTree)              // 获取用户动态路由
		menuRouter.POST("authority/get", authorityMenuApi.GetMenuAuthority)    // 获取指定角色menu
		menuRouter.POST("get/id", authorityMenuApi.GetBaseMenuById)            // 根据id获取菜单
	}
	return menuRouter
}
