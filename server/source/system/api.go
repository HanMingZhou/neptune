package system

import (
	"context"
	"gin-vue-admin/global"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/service/system"

	"github.com/pkg/errors"
)

type initApi struct{}

const initOrderApi = system.InitOrderSystem + 1

// auto run
func init() {
	system.RegisterInit(initOrderApi, &initApi{})
}

func (i *initApi) InitializerName() string {
	return sysModel.SysApi{}.TableName()
}

func (i *initApi) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, global.GVA_DB.AutoMigrate(&sysModel.SysApi{})
}

func (i *initApi) TableCreated(ctx context.Context) bool {
	return global.GVA_DB.Migrator().HasTable(&sysModel.SysApi{})
}

func (i *initApi) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []sysModel.SysApi{
		{ApiGroup: "jwt", Method: "POST", Path: "/api/v1/jwt/blacklist", Description: "jwt加入黑名单(退出，必选)"},

		{ApiGroup: "系统用户", Method: "POST", Path: "/api/v1/user/delete", Description: "删除用户"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/api/v1/user/register", Description: "用户注册"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/api/v1/user/list", Description: "获取用户列表"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/api/v1/user/info/update", Description: "设置用户信息"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/api/v1/user/self/info/update", Description: "设置自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "GET", Path: "/api/v1/user/info", Description: "获取自身信息(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/api/v1/user/authorities/update", Description: "设置权限组"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/api/v1/user/password/update", Description: "修改密码（建议选择)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/api/v1/user/authority/update", Description: "修改用户角色(必选)"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/api/v1/user/password/reset", Description: "重置用户密码"},
		{ApiGroup: "系统用户", Method: "POST", Path: "/api/v1/user/self/setting/update", Description: "用户界面配置"},

		{ApiGroup: "api", Method: "POST", Path: "/api/v1/api/add", Description: "创建api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/v1/api/delete", Description: "删除Api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/v1/api/update", Description: "更新Api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/v1/api/list", Description: "获取api列表"},
		{ApiGroup: "api", Method: "POST", Path: "/api/v1/api/all", Description: "获取所有api"},
		{ApiGroup: "api", Method: "POST", Path: "/api/v1/api/get", Description: "获取api详细信息"},
		{ApiGroup: "api", Method: "POST", Path: "/api/v1/api/delete/multi", Description: "批量删除api"},
		{ApiGroup: "api", Method: "GET", Path: "/api/v1/api/sync", Description: "获取待同步API"},
		{ApiGroup: "api", Method: "GET", Path: "/api/v1/api/group/list", Description: "获取路由组"},
		{ApiGroup: "api", Method: "POST", Path: "/api/v1/api/sync/enter", Description: "确认同步API"},
		{ApiGroup: "api", Method: "POST", Path: "/api/v1/api/ignore", Description: "忽略API"},
		{ApiGroup: "api", Method: "GET", Path: "/api/v1/api/casbin/fresh", Description: "刷新casbin权限"},

		{ApiGroup: "角色", Method: "POST", Path: "/api/v1/authority/copy", Description: "拷贝角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/api/v1/authority/add", Description: "创建角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/api/v1/authority/delete", Description: "删除角色"},
		{ApiGroup: "角色", Method: "POST", Path: "/api/v1/authority/update", Description: "更新角色信息"},
		{ApiGroup: "角色", Method: "POST", Path: "/api/v1/authority/list", Description: "获取角色列表"},
		{ApiGroup: "角色", Method: "POST", Path: "/api/v1/authority/data/authority/update", Description: "设置角色资源权限"},

		{ApiGroup: "casbin", Method: "POST", Path: "/api/v1/casbin/update", Description: "更改角色api权限"},
		{ApiGroup: "casbin", Method: "POST", Path: "/api/v1/casbin/policy/list", Description: "获取权限列表"},

		{ApiGroup: "菜单", Method: "POST", Path: "/api/v1/menu/add", Description: "新增菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/api/v1/menu/get", Description: "获取菜单树(必选)"},
		{ApiGroup: "菜单", Method: "POST", Path: "/api/v1/menu/delete", Description: "删除菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/api/v1/menu/update", Description: "更新菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/api/v1/menu/get/id", Description: "根据id获取菜单"},
		{ApiGroup: "菜单", Method: "POST", Path: "/api/v1/menu/list", Description: "分页获取基础menu列表"},
		{ApiGroup: "菜单", Method: "POST", Path: "/api/v1/menu/tree", Description: "获取用户动态路由"},
		{ApiGroup: "菜单", Method: "POST", Path: "/api/v1/menu/authority/get", Description: "获取指定角色menu"},
		{ApiGroup: "菜单", Method: "POST", Path: "/api/v1/menu/authority/update", Description: "增加menu和角色关联关系"},

		{ApiGroup: "分片上传", Method: "GET", Path: "/api/v1/file/find", Description: "查询当前文件成功的切片"},
		{ApiGroup: "分片上传", Method: "POST", Path: "/api/v1/file/breakpoint/continue", Description: "断点续传"},
		{ApiGroup: "分片上传", Method: "POST", Path: "/api/v1/file/breakpoint/finish", Description: "切片传输完成"},
		{ApiGroup: "分片上传", Method: "POST", Path: "/api/v1/file/chunk/delete", Description: "删除切片"},

		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/api/v1/file/upload", Description: "文件上传（建议选择）"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/api/v1/file/delete", Description: "删除文件"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/api/v1/file/update", Description: "文件名或者备注编辑"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/api/v1/file/list", Description: "获取上传文件列表"},
		{ApiGroup: "文件上传与下载", Method: "POST", Path: "/api/v1/file/import", Description: "导入URL"},

		{ApiGroup: "系统服务", Method: "POST", Path: "/api/v1/system/info", Description: "获取服务器信息"},
		{ApiGroup: "系统服务", Method: "POST", Path: "/api/v1/system/config/get", Description: "获取配置文件内容"},
		{ApiGroup: "系统服务", Method: "POST", Path: "/api/v1/system/config/update", Description: "设置配置文件内容"},
		{ApiGroup: "系统服务", Method: "POST", Path: "/api/v1/system/reload", Description: "重启服务"},
		{ApiGroup: "系统服务", Method: "POST", Path: "/api/v1/init/checkdb", Description: "检测是否需要初始化数据库"},

		{ApiGroup: "客户", Method: "POST", Path: "/api/v1/customer/add", Description: "创建客户"},
		{ApiGroup: "客户", Method: "POST", Path: "/api/v1/customer/update", Description: "更新客户"},
		{ApiGroup: "客户", Method: "POST", Path: "/api/v1/customer/delete", Description: "删除客户"},
		{ApiGroup: "客户", Method: "GET", Path: "/api/v1/customer/get", Description: "获取单一客户"},
		{ApiGroup: "客户", Method: "GET", Path: "/api/v1/customer/list", Description: "获取客户列表"},

		{ApiGroup: "代码生成器", Method: "GET", Path: "/api/v1/autocode/db/list", Description: "获取所有数据库"},
		{ApiGroup: "代码生成器", Method: "GET", Path: "/api/v1/autocode/table/list", Description: "获取数据库表"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/api/v1/autocode/add", Description: "自动化代码"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/api/v1/autocode/preview", Description: "预览自动化代码"},
		{ApiGroup: "代码生成器", Method: "GET", Path: "/api/v1/autocode/column/list", Description: "获取所选table的所有字段"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/api/v1/autocode/plugin/install", Description: "安装插件"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/api/v1/autocode/plugin/pack", Description: "打包插件"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/api/v1/autocode/mcp/add", Description: "自动生成 MCP Tool 模板"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/api/v1/autocode/mcp/test", Description: "MCP Tool 测试"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/api/v1/autocode/mcp/list", Description: "获取 MCP ToolList"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/api/v1/autocode/llm/add", Description: "LLM生成代码"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/api/v1/autocode/menu/init", Description: "同步插件菜单"},
		{ApiGroup: "代码生成器", Method: "POST", Path: "/api/v1/autocode/api/init", Description: "同步插件API"},

		{ApiGroup: "模板配置", Method: "POST", Path: "/api/v1/autocode/package/add", Description: "配置模板"},
		{ApiGroup: "模板配置", Method: "GET", Path: "/api/v1/autocode/template/list", Description: "获取模板文件"},
		{ApiGroup: "模板配置", Method: "POST", Path: "/api/v1/autocode/package/list", Description: "获取所有模板"},
		{ApiGroup: "模板配置", Method: "POST", Path: "/api/v1/autocode/package/delete", Description: "删除模板"},

		{ApiGroup: "代码生成器历史", Method: "POST", Path: "/api/v1/autocode/meta/get", Description: "获取meta信息"},
		{ApiGroup: "代码生成器历史", Method: "POST", Path: "/api/v1/autocode/rollback", Description: "回滚自动生成代码"},
		{ApiGroup: "代码生成器历史", Method: "POST", Path: "/api/v1/autocode/list", Description: "查询回滚记录"},
		{ApiGroup: "代码生成器历史", Method: "POST", Path: "/api/v1/autocode/delete", Description: "删除回滚记录"},
		{ApiGroup: "代码生成器历史", Method: "POST", Path: "/api/v1/autocode/func/add", Description: "增加模板方法"},

		{ApiGroup: "系统字典详情", Method: "POST", Path: "/api/v1/dictionary/detail/update", Description: "更新字典内容"},
		{ApiGroup: "系统字典详情", Method: "POST", Path: "/api/v1/dictionary/detail/add", Description: "新增字典内容"},
		{ApiGroup: "系统字典详情", Method: "POST", Path: "/api/v1/dictionary/detail/delete", Description: "删除字典内容"},
		{ApiGroup: "系统字典详情", Method: "GET", Path: "/api/v1/dictionary/detail/get", Description: "根据ID获取字典内容"},
		{ApiGroup: "系统字典详情", Method: "GET", Path: "/api/v1/dictionary/detail/list", Description: "获取字典内容列表"},

		{ApiGroup: "系统字典详情", Method: "GET", Path: "/api/v1/dictionary/detail/tree/list", Description: "获取字典数列表"},
		{ApiGroup: "系统字典详情", Method: "GET", Path: "/api/v1/dictionary/detail/tree/type/list", Description: "根据分类获取字典数列表"},
		{ApiGroup: "系统字典详情", Method: "GET", Path: "/api/v1/dictionary/detail/parent/get", Description: "根据父级ID获取字典详情"},
		{ApiGroup: "系统字典详情", Method: "GET", Path: "/api/v1/dictionary/detail/path/get", Description: "获取字典详情的完整路径"},

		{ApiGroup: "系统字典", Method: "POST", Path: "/api/v1/dictionary/add", Description: "新增字典"},
		{ApiGroup: "系统字典", Method: "POST", Path: "/api/v1/dictionary/delete", Description: "删除字典"},
		{ApiGroup: "系统字典", Method: "POST", Path: "/api/v1/dictionary/update", Description: "更新字典"},
		{ApiGroup: "系统字典", Method: "GET", Path: "/api/v1/dictionary/get", Description: "根据ID获取字典（建议选择）"},
		{ApiGroup: "系统字典", Method: "GET", Path: "/api/v1/dictionary/list", Description: "获取字典列表"},
		{ApiGroup: "系统字典", Method: "POST", Path: "/api/v1/dictionary/import", Description: "导入字典JSON"},
		{ApiGroup: "系统字典", Method: "GET", Path: "/api/v1/dictionary/export", Description: "导出字典JSON"},

		{ApiGroup: "操作记录", Method: "POST", Path: "/api/v1/operation/record/add", Description: "新增操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/api/v1/operation/record/get", Description: "根据ID获取操作记录"},
		{ApiGroup: "操作记录", Method: "GET", Path: "/api/v1/operation/record/list", Description: "获取操作记录列表"},
		{ApiGroup: "操作记录", Method: "POST", Path: "/api/v1/operation/record/delete", Description: "删除操作记录"},
		{ApiGroup: "操作记录", Method: "POST", Path: "/api/v1/operation/record/delete/multi", Description: "批量删除操作历史"},

		{ApiGroup: "断点续传(插件版)", Method: "POST", Path: "/api/v1/file/simple/upload", Description: "插件版分片上传"},
		{ApiGroup: "断点续传(插件版)", Method: "GET", Path: "/api/v1/file/simple/check", Description: "文件完整度验证"},
		{ApiGroup: "断点续传(插件版)", Method: "GET", Path: "/api/v1/file/simple/merge", Description: "上传完成合并文件"},

		{ApiGroup: "email", Method: "POST", Path: "/api/v1/email/test", Description: "发送测试邮件"},
		{ApiGroup: "email", Method: "POST", Path: "/api/v1/email/send", Description: "发送邮件"},

		{ApiGroup: "按钮权限", Method: "POST", Path: "/api/v1/authority/btn/update", Description: "设置按钮权限"},
		{ApiGroup: "按钮权限", Method: "POST", Path: "/api/v1/authority/btn/get", Description: "获取已有按钮权限"},
		{ApiGroup: "按钮权限", Method: "POST", Path: "/api/v1/authority/btn/delete", Description: "删除按钮"},

		{ApiGroup: "导出模板", Method: "POST", Path: "/api/v1/export/template/add", Description: "新增导出模板"},
		{ApiGroup: "导出模板", Method: "POST", Path: "/api/v1/export/template/delete", Description: "删除导出模板"},
		{ApiGroup: "导出模板", Method: "POST", Path: "/api/v1/export/template/delete/multi", Description: "批量删除导出模板"},
		{ApiGroup: "导出模板", Method: "POST", Path: "/api/v1/export/template/update", Description: "更新导出模板"},
		{ApiGroup: "导出模板", Method: "GET", Path: "/api/v1/export/template/get", Description: "根据ID获取导出模板"},
		{ApiGroup: "导出模板", Method: "GET", Path: "/api/v1/export/template/list", Description: "获取导出模板列表"},
		{ApiGroup: "导出模板", Method: "GET", Path: "/api/v1/export/template/export", Description: "导出Excel"},
		{ApiGroup: "导出模板", Method: "GET", Path: "/api/v1/export/template/download", Description: "下载模板"},
		{ApiGroup: "导出模板", Method: "GET", Path: "/api/v1/export/template/sql/preview", Description: "预览SQL"},
		{ApiGroup: "导出模板", Method: "POST", Path: "/api/v1/export/template/import", Description: "导入Excel"},

		{ApiGroup: "错误日志", Method: "POST", Path: "/api/v1/error/add", Description: "新建错误日志"},
		{ApiGroup: "错误日志", Method: "POST", Path: "/api/v1/error/delete", Description: "删除错误日志"},
		{ApiGroup: "错误日志", Method: "POST", Path: "/api/v1/error/delete/multi", Description: "批量删除错误日志"},
		{ApiGroup: "错误日志", Method: "POST", Path: "/api/v1/error/update", Description: "更新错误日志"},
		{ApiGroup: "错误日志", Method: "GET", Path: "/api/v1/error/get", Description: "根据ID获取错误日志"},
		{ApiGroup: "错误日志", Method: "GET", Path: "/api/v1/error/list", Description: "获取错误日志列表"},
		{ApiGroup: "错误日志", Method: "GET", Path: "/api/v1/error/solution/get", Description: "触发错误处理(异步)"},

		{ApiGroup: "公告", Method: "POST", Path: "/api/v1/info/add", Description: "新建公告"},
		{ApiGroup: "公告", Method: "POST", Path: "/api/v1/info/delete", Description: "删除公告"},
		{ApiGroup: "公告", Method: "POST", Path: "/api/v1/info/delete/multi", Description: "批量删除公告"},
		{ApiGroup: "公告", Method: "POST", Path: "/api/v1/info/update", Description: "更新公告"},
		{ApiGroup: "公告", Method: "GET", Path: "/api/v1/info/get", Description: "根据ID获取公告"},
		{ApiGroup: "公告", Method: "GET", Path: "/api/v1/info/list", Description: "获取公告列表"},

		{ApiGroup: "参数管理", Method: "POST", Path: "/api/v1/params/add", Description: "新建参数"},
		{ApiGroup: "参数管理", Method: "POST", Path: "/api/v1/params/delete", Description: "删除参数"},
		{ApiGroup: "参数管理", Method: "POST", Path: "/api/v1/params/delete/multi", Description: "批量删除参数"},
		{ApiGroup: "参数管理", Method: "POST", Path: "/api/v1/params/update", Description: "更新参数"},
		{ApiGroup: "参数管理", Method: "GET", Path: "/api/v1/params/get", Description: "根据ID获取参数"},
		{ApiGroup: "参数管理", Method: "GET", Path: "/api/v1/params/list", Description: "获取参数列表"},
		{ApiGroup: "参数管理", Method: "GET", Path: "/api/v1/params/get/key", Description: "获取参数列表"},
		{ApiGroup: "媒体库分类", Method: "GET", Path: "/api/v1/attachment/category/list", Description: "分类列表"},
		{ApiGroup: "媒体库分类", Method: "POST", Path: "/api/v1/attachment/category/add", Description: "添加/编辑分类"},
		{ApiGroup: "媒体库分类", Method: "POST", Path: "/api/v1/attachment/category/delete", Description: "删除分类"},

		{ApiGroup: "版本控制", Method: "GET", Path: "/api/v1/version/get", Description: "获取单一版本"},
		{ApiGroup: "版本控制", Method: "GET", Path: "/api/v1/version/list", Description: "获取版本列表"},
		{ApiGroup: "版本控制", Method: "GET", Path: "/api/v1/version/download/json", Description: "下载版本json"},
		{ApiGroup: "版本控制", Method: "GET", Path: "/api/v1/version/downloadVersionJson", Description: "下载版本JSON数据"},
		{ApiGroup: "版本控制", Method: "POST", Path: "/api/v1/version/export", Description: "创建版本"},
		{ApiGroup: "版本控制", Method: "POST", Path: "/api/v1/version/import", Description: "同步版本"},
		{ApiGroup: "版本控制", Method: "POST", Path: "/api/v1/version/delete", Description: "删除版本"},
		{ApiGroup: "版本控制", Method: "POST", Path: "/api/v1/version/delete/multi", Description: "批量删除版本"},

		// ============ 业务模块 API ============
		{ApiGroup: "个人中心", Method: "GET", Path: "/api/v1/account/security/status", Description: "获取安全状态"},
		{ApiGroup: "个人中心", Method: "POST", Path: "/api/v1/account/access/log/list", Description: "获取访问日志列表"},
		{ApiGroup: "个人中心", Method: "POST", Path: "/api/v1/account/active/session/list", Description: "获取活跃会话列表"},
		{ApiGroup: "个人中心", Method: "POST", Path: "/api/v1/account/password/update", Description: "修改密码"},
		{ApiGroup: "个人中心", Method: "POST", Path: "/api/v1/account/bind", Description: "绑定手机/邮箱"},
		{ApiGroup: "个人中心", Method: "POST", Path: "/api/v1/account/mfa/setup", Description: "MFA初始化"},
		{ApiGroup: "个人中心", Method: "POST", Path: "/api/v1/account/mfa/activate", Description: "MFA激活"},
		{ApiGroup: "个人中心", Method: "POST", Path: "/api/v1/account/mfa/toggle", Description: "开启/关闭MFA"},
		{ApiGroup: "个人中心", Method: "POST", Path: "/api/v1/account/ak/generate", Description: "生成AccessKey"},
		{ApiGroup: "个人中心", Method: "POST", Path: "/api/v1/account/active/session/kill", Description: "强制下线"},

		{ApiGroup: "费用账单", Method: "GET", Path: "/api/v1/order/overview", Description: "获取财务总览"},
		{ApiGroup: "费用账单", Method: "POST", Path: "/api/v1/order/usage/list", Description: "使用详情列表"},
		{ApiGroup: "费用账单", Method: "POST", Path: "/api/v1/order/transaction/list", Description: "收支明细列表"},
		{ApiGroup: "费用账单", Method: "POST", Path: "/api/v1/order/order/list", Description: "订单列表"},
		{ApiGroup: "费用账单", Method: "POST", Path: "/api/v1/order/invoice/list", Description: "发票记录列表"},
		{ApiGroup: "费用账单", Method: "POST", Path: "/api/v1/order/invoice/apply", Description: "申请发票"},

		{ApiGroup: "容器实例", Method: "POST", Path: "/api/v1/notebook/add", Description: "创建容器实例"},
		{ApiGroup: "容器实例", Method: "POST", Path: "/api/v1/notebook/delete", Description: "删除容器实例"},
		{ApiGroup: "容器实例", Method: "GET", Path: "/api/v1/notebook/list", Description: "获取容器实例列表"},
		{ApiGroup: "容器实例", Method: "POST", Path: "/api/v1/notebook/update", Description: "更新容器实例"},
		{ApiGroup: "容器实例", Method: "POST", Path: "/api/v1/notebook/stop", Description: "停止容器实例"},
		{ApiGroup: "容器实例", Method: "POST", Path: "/api/v1/notebook/start", Description: "启动容器实例"},
		{ApiGroup: "容器实例", Method: "GET", Path: "/api/v1/notebook/get", Description: "获取容器实例详情"},
		{ApiGroup: "容器实例", Method: "GET", Path: "/api/v1/notebook/pod/list", Description: "获取容器实例 Pod 列表"},
		{ApiGroup: "容器实例", Method: "GET", Path: "/api/v1/notebook/log/list", Description: "获取容器实例日志"},
		{ApiGroup: "容器实例", Method: "GET", Path: "/api/v1/notebook/terminal", Description: "Web Terminal"},
		{ApiGroup: "容器实例", Method: "GET", Path: "/api/v1/notebook/log/stream", Description: "日志流"},
		{ApiGroup: "终端及身份认证", Method: "GET", Path: "/api/v1/notebook/auth", Description: "交互式终端鉴权"},

		{ApiGroup: "训练任务", Method: "POST", Path: "/api/v1/training/add", Description: "创建训练任务"},
		{ApiGroup: "训练任务", Method: "POST", Path: "/api/v1/training/delete", Description: "删除训练任务"},
		{ApiGroup: "训练任务", Method: "POST", Path: "/api/v1/training/stop", Description: "停止训练任务"},
		{ApiGroup: "训练任务", Method: "GET", Path: "/api/v1/training/list", Description: "获取训练任务列表"},
		{ApiGroup: "训练任务", Method: "GET", Path: "/api/v1/training/get", Description: "获取训练任务详情"},
		{ApiGroup: "训练任务", Method: "GET", Path: "/api/v1/training/pod/list", Description: "获取训练任务 Pod 列表"},
		{ApiGroup: "训练任务", Method: "GET", Path: "/api/v1/training/log/list", Description: "获取训练任务日志"},
		{ApiGroup: "训练任务", Method: "GET", Path: "/api/v1/training/log/download", Description: "下载训练任务日志"},
		{ApiGroup: "训练任务", Method: "GET", Path: "/api/v1/training/terminal", Description: "Web Terminal"},
		{ApiGroup: "训练任务", Method: "GET", Path: "/api/v1/training/log/stream", Description: "日志流"},

		{ApiGroup: "镜像", Method: "GET", Path: "/api/v1/image/list", Description: "获取镜像列表"},
		{ApiGroup: "镜像", Method: "POST", Path: "/api/v1/image/add", Description: "创建镜像"},
		{ApiGroup: "镜像", Method: "POST", Path: "/api/v1/image/update", Description: "更新镜像"},
		{ApiGroup: "镜像", Method: "POST", Path: "/api/v1/image/delete", Description: "删除镜像"},

		{ApiGroup: "产品", Method: "GET", Path: "/api/v1/product/list", Description: "获取产品列表"},
		{ApiGroup: "产品", Method: "GET", Path: "/api/v1/product/get", Description: "获取产品详情"},
		{ApiGroup: "产品", Method: "GET", Path: "/api/v1/product/filter/list", Description: "获取产品筛选条件"},

		{ApiGroup: "SSH密钥", Method: "POST", Path: "/api/v1/sshkey/add", Description: "添加SSH密钥"},
		{ApiGroup: "SSH密钥", Method: "POST", Path: "/api/v1/sshkey/delete", Description: "删除SSH密钥"},
		{ApiGroup: "SSH密钥", Method: "POST", Path: "/api/v1/sshkey/default/update", Description: "设置默认密钥"},
		{ApiGroup: "SSH密钥", Method: "POST", Path: "/api/v1/sshkey/list", Description: "获取SSH密钥列表"},

		{ApiGroup: "存储卷", Method: "POST", Path: "/api/v1/pvc/add", Description: "创建存储卷"},
		{ApiGroup: "存储卷", Method: "POST", Path: "/api/v1/pvc/delete", Description: "删除存储卷"},
		{ApiGroup: "存储卷", Method: "GET", Path: "/api/v1/pvc/list", Description: "获取存储卷列表"},
		{ApiGroup: "存储卷", Method: "POST", Path: "/api/v1/pvc/update", Description: "更新存储卷"},

		{ApiGroup: "文件存储", Method: "POST", Path: "/api/v1/volume/add", Description: "创建文件存储"},
		{ApiGroup: "文件存储", Method: "GET", Path: "/api/v1/volume/list", Description: "获取文件存储列表"},
		{ApiGroup: "文件存储", Method: "POST", Path: "/api/v1/volume/expand", Description: "扩容文件存储"},
		{ApiGroup: "文件存储", Method: "POST", Path: "/api/v1/volume/delete", Description: "删除文件存储"},
		{ApiGroup: "文件存储", Method: "GET", Path: "/api/v1/volume/area/list", Description: "获取可用区域列表"},

		{ApiGroup: "SSH代理", Method: "POST", Path: "/api/v1/piper/add", Description: "创建SSH代理"},
		{ApiGroup: "SSH代理", Method: "POST", Path: "/api/v1/piper/delete", Description: "删除SSH代理"},
		{ApiGroup: "SSH代理", Method: "GET", Path: "/api/v1/piper/list", Description: "获取SSH代理列表"},
		{ApiGroup: "SSH代理", Method: "POST", Path: "/api/v1/piper/update", Description: "更新SSH代理"},

		{ApiGroup: "Secret", Method: "POST", Path: "/api/v1/secret/add", Description: "创建Secret"},
		{ApiGroup: "Secret", Method: "POST", Path: "/api/v1/secret/delete", Description: "删除Secret"},
		{ApiGroup: "Secret", Method: "GET", Path: "/api/v1/secret/list", Description: "获取Secret列表"},
		{ApiGroup: "Secret", Method: "POST", Path: "/api/v1/secret/update", Description: "更新Secret"},

		{ApiGroup: "TensorBoard", Method: "POST", Path: "/api/v1/tensorboard/add", Description: "创建TensorBoard"},
		{ApiGroup: "TensorBoard", Method: "POST", Path: "/api/v1/tensorboard/delete", Description: "删除TensorBoard"},
		{ApiGroup: "TensorBoard", Method: "GET", Path: "/api/v1/tensorboard/list", Description: "获取TensorBoard列表"},
		{ApiGroup: "TensorBoard", Method: "POST", Path: "/api/v1/tensorboard/update", Description: "更新TensorBoard"},

		{ApiGroup: "推理服务", Method: "POST", Path: "/api/v1/inference/add", Description: "创建推理服务"},
		{ApiGroup: "推理服务", Method: "POST", Path: "/api/v1/inference/delete", Description: "删除推理服务"},
		{ApiGroup: "推理服务", Method: "POST", Path: "/api/v1/inference/stop", Description: "停止推理服务"},
		{ApiGroup: "推理服务", Method: "POST", Path: "/api/v1/inference/start", Description: "启动推理服务"},
		{ApiGroup: "推理服务", Method: "POST", Path: "/api/v1/inference/list", Description: "获取推理服务列表"},
		{ApiGroup: "推理服务", Method: "GET", Path: "/api/v1/inference/get", Description: "获取推理服务详情"},
		{ApiGroup: "推理服务", Method: "GET", Path: "/api/v1/inference/pod/list", Description: "获取推理服务 Pod 列表"},
		{ApiGroup: "推理服务", Method: "GET", Path: "/api/v1/inference/log/list", Description: "获取推理服务日志"},
		{ApiGroup: "推理服务", Method: "POST", Path: "/api/v1/inference/api/key/add", Description: "创建 API Key"},
		{ApiGroup: "推理服务", Method: "POST", Path: "/api/v1/inference/api/key/delete", Description: "删除 API Key"},
		{ApiGroup: "推理服务", Method: "POST", Path: "/api/v1/inference/api/key/list", Description: "获取 API Key 列表"},
		{ApiGroup: "推理服务", Method: "GET", Path: "/api/v1/inference/terminal", Description: "推理服务终端"},
		{ApiGroup: "推理服务", Method: "GET", Path: "/api/v1/inference/log/stream", Description: "推理服务日志流"},

		{ApiGroup: "控制台仪表盘", Method: "POST", Path: "/api/v1/dashboard/get", Description: "获取仪表盘数据"},

		{ApiGroup: "CMS管理", Method: "POST", Path: "/api/v1/cms/product/add", Description: "创建产品"},
		{ApiGroup: "CMS管理", Method: "POST", Path: "/api/v1/cms/product/update", Description: "更新产品"},
		{ApiGroup: "CMS管理", Method: "POST", Path: "/api/v1/cms/product/price/update", Description: "更新价格"},
		{ApiGroup: "CMS管理", Method: "POST", Path: "/api/v1/cms/product/delete", Description: "删除产品"},
		{ApiGroup: "CMS管理", Method: "GET", Path: "/api/v1/cms/product/list", Description: "获取产品列表"},
		{ApiGroup: "CMS管理", Method: "GET", Path: "/api/v1/cms/product/detail", Description: "获取产品详情"},
		{ApiGroup: "CMS管理", Method: "GET", Path: "/api/v1/cms/product/cluster/list", Description: "获取集群列表"},
		{ApiGroup: "CMS管理", Method: "GET", Path: "/api/v1/cms/product/area/list", Description: "获取地区列表"},
		{ApiGroup: "CMS管理", Method: "GET", Path: "/api/v1/cms/product/node/list", Description: "获取集群节点列表"},
		{ApiGroup: "CMS管理", Method: "POST", Path: "/api/v1/cms/node/list", Description: "获取节点列表"},
		{ApiGroup: "CMS管理", Method: "POST", Path: "/api/v1/cms/node/uncordon", Description: "恢复节点调度"},
		{ApiGroup: "CMS管理", Method: "POST", Path: "/api/v1/cms/node/drain", Description: "驱逐节点"},

		{ApiGroup: "集群管理", Method: "POST", Path: "/api/v1/cms/cluster/list", Description: "获取集群列表"},
		{ApiGroup: "集群管理", Method: "POST", Path: "/api/v1/cms/cluster/add", Description: "创建集群"},
		{ApiGroup: "集群管理", Method: "POST", Path: "/api/v1/cms/cluster/update", Description: "更新集群"},
		{ApiGroup: "集群管理", Method: "POST", Path: "/api/v1/cms/cluster/delete", Description: "删除集群"},
	}

	// 批量查询已存在的 API，避免循环执行 SQL
	var existingApis []sysModel.SysApi
	if err := global.GVA_DB.Select("path", "method").Find(&existingApis).Error; err != nil {
		return ctx, errors.Wrap(err, "查询已存在API失败")
	}

	existingMap := make(map[string]bool)
	for _, a := range existingApis {
		existingMap[a.Path+":"+a.Method] = true
	}

	var missingEntities []sysModel.SysApi
	for _, entity := range entities {
		if !existingMap[entity.Path+":"+entity.Method] {
			missingEntities = append(missingEntities, entity)
		}
	}

	if len(missingEntities) > 0 {
		if err := global.GVA_DB.Create(&missingEntities).Error; err != nil {
			return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"表数据批量初始化失败!")
		}
	}

	for _, entity := range entities {
		// 检查描述是否一致，不一致则更新
		if existingMap[entity.Path+":"+entity.Method] {
			global.GVA_DB.Model(&sysModel.SysApi{}).Where("path = ? AND method = ?", entity.Path, entity.Method).Updates(map[string]interface{}{
				"description": entity.Description,
				"api_group":   entity.ApiGroup,
			})
		}
	}

	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
