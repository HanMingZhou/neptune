package initialize

import (
	"gin-vue-admin/global"
	"gin-vue-admin/middleware"
	"gin-vue-admin/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 初始化总路由

func Routers() *gin.Engine {
	Router := gin.New()
	// 使用自定义的 Recovery 中间件，记录 panic 并入库
	// Router.Use(middleware.GinRecovery(true))
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}

	// 注册mcp服务
	if !global.GVA_CONFIG.MCP.Separate {
		sseServer := McpRun()
		Router.GET(global.GVA_CONFIG.MCP.SSEPath, func(c *gin.Context) {
			sseServer.SSEHandler().ServeHTTP(c.Writer, c.Request)
		})
		Router.POST(global.GVA_CONFIG.MCP.MessagePath, func(c *gin.Context) {
			sseServer.MessageHandler().ServeHTTP(c.Writer, c.Request)
		})
	}
	exampleRouter := router.RouterGroupApp.Example
	systemRouter := router.RouterGroupApp.System

	imageRouter := router.RouterGroupApp.Image
	notebookRouter := router.RouterGroupApp.NoteBook
	apisixRouter := router.RouterGroupApp.Apisix
	piperRouter := router.RouterGroupApp.Piper
	productRouter := router.RouterGroupApp.Product
	pvcRouter := router.RouterGroupApp.PVC
	secretRouter := router.RouterGroupApp.Secret
	tensorboardRouter := router.RouterGroupApp.TensorBoard
	clusterRouter := router.RouterGroupApp.Cluster
	cmsRouter := router.RouterGroupApp.CMS
	sshKeyRouter := router.RouterGroupApp.SSHKey
	trainingRouter := router.RouterGroupApp.Training
	inferenceRouter := router.RouterGroupApp.Inference
	dashboardRouter := router.RouterGroupApp.Dashboard
	accountRouter := router.RouterGroupApp.Account
	orderRouter := router.RouterGroupApp.Order

	// 如果想要不使用nginx代理前端网页，可以修改 web/.env.production 下的
	// VUE_APP_BASE_API = /
	// VUE_APP_BASE_PATH = http://localhost
	// 然后执行打包命令 npm run build。在打开下面3行注释
	// Router.StaticFile("/favicon.ico", "./dist/favicon.ico")
	// Router.Static("/assets", "./dist/assets")   // dist里面的静态资源
	// Router.StaticFile("/", "./dist/index.html") // 前端网页入口页面

	// Router.StaticFS(global.GVA_CONFIG.Local.StorePath, justFilesFilesystem{http.Dir(global.GVA_CONFIG.Local.StorePath)}) // Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")
	// 跨域，如需跨域可以打开下面的注释
	Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	// global.GVA_LOG.Info("use middleware cors")
	// docs.SwaggerInfo.BasePath = global.GVA_CONFIG.System.RouterPrefix
	// Router.GET(global.GVA_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// global.GVA_LOG.Info("register swagger handler")
	// 方便统一添加路由组前缀 多服务器上线使用

	// 公共路由
	PublicGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)

	// 私有路由
	PrivateGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())

	// API v1 路由组 (统一版本管理，方便将来扩展 v2)
	v1PrivateGroup := PrivateGroup.Group("api/v1")
	v1PublicGroup := PublicGroup.Group("api/v1")

	// 健康监测
	{

		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}

	// 系统管理 (通用功能)
	{
		systemRouter.InitBaseRouter(v1PublicGroup) // 注册基础功能路由 不做鉴权
	}

	// 系统内置 (v1 规范)
	{
		systemRouter.InitApiRouter(v1PrivateGroup, v1PublicGroup)               // 注册功能api路由
		systemRouter.InitJwtRouter(v1PrivateGroup)                              // jwt相关路由
		systemRouter.InitUserRouter(v1PrivateGroup, v1PublicGroup)              // 注册用户路由
		systemRouter.InitMenuRouter(v1PrivateGroup)                             // 注册menu路由
		systemRouter.InitSystemRouter(v1PrivateGroup)                           // system相关路由
		systemRouter.InitSysVersionRouter(v1PrivateGroup)                       // 发版相关路由
		systemRouter.InitCasbinRouter(v1PrivateGroup)                           // 权限相关路由
		systemRouter.InitAutoCodeRouter(v1PrivateGroup, v1PublicGroup)          // 创建自动化代码
		systemRouter.InitAuthorityRouter(v1PrivateGroup)                        // 注册角色路由
		systemRouter.InitSysDictionaryRouter(v1PrivateGroup)                    // 字典管理
		systemRouter.InitAutoCodeHistoryRouter(v1PrivateGroup)                  // 自动化代码历史
		systemRouter.InitSysOperationRecordRouter(v1PrivateGroup)               // 操作记录
		systemRouter.InitSysDictionaryDetailRouter(v1PrivateGroup)              // 字典详情管理
		systemRouter.InitAuthorityBtnRouterRouter(v1PrivateGroup)               // 按钮权限管理
		systemRouter.InitSysExportTemplateRouter(v1PrivateGroup, v1PublicGroup) // 导出模板
		systemRouter.InitSysParamsRouter(v1PrivateGroup, v1PublicGroup)         // 参数管理
		systemRouter.InitSysErrorRouter(v1PrivateGroup, v1PublicGroup)          // 错误日志
		exampleRouter.InitCustomerRouter(v1PrivateGroup)                        // 客户路由
		exampleRouter.InitFileUploadAndDownloadRouter(v1PrivateGroup)           // 文件上传下载功能路由
		exampleRouter.InitAttachmentCategoryRouterRouter(v1PrivateGroup)        // 文件上传下载分类
	}

	// 业务模块路由 (v1)
	clusterRouter.InitClusterRouter(v1PrivateGroup)
	cmsRouter.InitCMSProductRouter(v1PrivateGroup)
	cmsRouter.InitNodeRouter(v1PrivateGroup)
	notebookRouter.InitNotebookApiRouter(v1PrivateGroup, v1PublicGroup) // notebook 接口
	apisixRouter.InitApisixRouter(v1PublicGroup)                        // Apisix forward-auth 认证接口
	imageRouter.InitImageRouter(v1PrivateGroup)
	piperRouter.InitPipeRouter(v1PrivateGroup)
	productRouter.InitProductRouter(v1PrivateGroup)
	pvcRouter.InitPVCRouter(v1PrivateGroup)
	pvcRouter.InitVolumeRouter(v1PrivateGroup)
	secretRouter.InitSecretRouter(v1PrivateGroup)
	tensorboardRouter.InitTensorBoardApiRouter(v1PrivateGroup)
	sshKeyRouter.InitSSHKeyRouter(v1PrivateGroup)
	trainingRouter.InitTrainingRouter(v1PrivateGroup, v1PublicGroup)
	inferenceRouter.InitInferenceRouter(v1PrivateGroup, v1PublicGroup)
	dashboardRouter.InitDashboardRouter(v1PrivateGroup)
	accountRouter.InitAccountRouter(v1PrivateGroup)
	orderRouter.InitOrderRouter(v1PrivateGroup)

	// 注册业务路由
	initBizRouter(PrivateGroup, PublicGroup)

	global.GVA_ROUTERS = Router.Routes()
	return Router
}
