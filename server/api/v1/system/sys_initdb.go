package system

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/system"
	"gin-vue-admin/model/system/request"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

type DBApi struct{}

// InitDB
// @Tags     InitDB
// @Summary  初始化用户数据库
// @Produce  application/json
// @Param    data  body      request.InitDB                  true  "初始化数据库参数"
// @Success  200   {object}  response.Response{data=string}  "初始化用户数据库"
// @Router   /init/initdb [post]
func (i *DBApi) InitDB(c *gin.Context) {
	var dbInfo request.InitDB
	if err := c.ShouldBindJSON(&dbInfo); err != nil {
		utils.HandleError(c, err, "参数校验不通过")
		return
	}

	// 如果数据库已初始化，尝试只初始化数据
	if global.GVA_DB != nil {
		// 调用 service 层的 InitDataOnly 方法
		if err := initDBService.InitDataOnly(); err != nil {
			utils.HandleError(c, err, "数据初始化失败")
			return
		}
		response.OkWithMessage("数据初始化成功", c)
		return
	}

	// 数据库未初始化，执行完整初始化流程
	if err := initDBService.InitDB(dbInfo); err != nil {
		utils.HandleError(c, err, "自动创建数据库失败")
		return
	}
	response.OkWithMessage("自动创建数据库成功", c)
}

// CheckDB
// @Tags     CheckDB
// @Summary  检测是否需要初始化数据库
// @Produce  application/json
// @Success  200  {object}  response.Response{data=map[string]interface{},msg=string}  "检测数据库状态"
// @Router   /init/checkdb [post]
func (i *DBApi) CheckDB(c *gin.Context) {
	var (
		message       = "前往初始化数据库"
		needInit      = true
		hasAdmin      = false
		dbInitialized = false
		hasBaseData   = false
	)

	if global.GVA_DB != nil {
		dbInitialized = true

		// 检查是否有基础数据（角色表）
		var authorityCount int64
		global.GVA_DB.Model(&system.SysAuthority{}).Count(&authorityCount)
		hasBaseData = authorityCount > 0

		if !hasBaseData {
			message = "数据库表已创建，但基础数据缺失，点击初始化"
			needInit = true
		} else {
			// 检查是否已有 admin 账号
			var adminUser system.SysUser
			if err := global.GVA_DB.Where("username = ?", "admin").First(&adminUser).Error; err == nil {
				hasAdmin = true
				message = "系统已初始化完成"
				needInit = false
			} else {
				message = "需要创建管理员账号"
				needInit = true
			}
		}
	}

	global.GVA_LOG.Info(message)
	response.OkWithDetailed(gin.H{
		"needInit":      needInit,
		"dbInitialized": dbInitialized,
		"hasAdmin":      hasAdmin,
		"hasBaseData":   hasBaseData,
	}, message, c)
}
