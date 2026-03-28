package system

import (
	"context"
	"gin-vue-admin/global"
	"gin-vue-admin/service/system"
	"gin-vue-admin/utils"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/pkg/errors"
)

const initOrderCasbin = initOrderApiIgnore + 1

type initCasbin struct{}

// auto run
func init() {
	system.RegisterInit(initOrderCasbin, &initCasbin{})
}

func (i *initCasbin) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, global.GVA_DB.AutoMigrate(&adapter.CasbinRule{})
}

func (i *initCasbin) TableCreated(ctx context.Context) bool {
	return global.GVA_DB.Migrator().HasTable(&adapter.CasbinRule{})
}

func (i *initCasbin) InitializerName() string {
	var entity adapter.CasbinRule
	return entity.TableName()
}

func (i *initCasbin) InitializeData(ctx context.Context) (context.Context, error) {
	entities := []adapter.CasbinRule{

		{Ptype: "p", V0: "888", V1: "/api/v1/user/register", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/get", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/all", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/delete/multi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/sync", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/group/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/sync/enter", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/ignore", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/api/casbin/fresh", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/api/v1/authority/copy", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/authority/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/authority/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/authority/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/authority/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/authority/data/authority/update", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/menu/get", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/menu/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/menu/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/menu/tree", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/menu/authority/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/menu/authority/get", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/menu/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/menu/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/menu/get/id", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/user/info", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/user/info/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/user/self/info/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/user/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/user/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/user/password/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/user/authority/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/user/authorities/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/user/password/reset", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/user/self/setting/update", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/file/upload/find", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/upload/finish", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/upload/breakpoint", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/upload/remove", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/file/upload", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/edit", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/import", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/find", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/breakpoint/continue", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/breakpoint/finish", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/chunk/delete", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/casbin/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/casbin/policy/list", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/jwt/blacklist", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/system/config/get", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/system/config/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/system/info", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/system/reload", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/customer/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/customer/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/customer/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/customer/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/customer/list", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/db/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/table/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/preview", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/column/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/plugin/install", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/plugin/pack", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/mcp/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/mcp/test", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/mcp/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/llm/add", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/package/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/template/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/package/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/package/delete", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/getMeta", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/meta/get", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/rollback", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/getSysHistory", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/delSysHistory", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/func/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/menu/init", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/autocode/api/init", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/detail/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/detail/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/detail/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/detail/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/detail/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/detail/tree/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/detail/tree/type/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/detail/parent/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/detail/path/get", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/import", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/dictionary/export", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/api/v1/operation/record/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/operation/record/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/operation/record/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/operation/record/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/operation/record/delete/multi", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/email/test", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/email/send", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/file/simple/upload", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/simple/check", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/file/simple/merge", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/api/v1/authority/btn/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/authority/btn/get", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/authority/btn/delete", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/export/template/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/export/template/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/export/template/delete/multi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/export/template/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/export/template/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/export/template/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/export/template/export", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/export/template/download", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/export/template/sql/preview", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/export/template/import", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/error/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/error/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/error/delete/multi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/error/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/error/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/error/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/error/solution/get", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/api/v1/info/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/info/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/info/delete/multi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/info/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/info/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/info/list", V2: "GET"},

		{Ptype: "p", V0: "888", V1: "/api/v1/params/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/params/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/params/delete/multi", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/params/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/params/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/params/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/params/get/key", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/attachment/category/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/attachment/category/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/attachment/category/delete", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/version/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/version/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/version/download/json", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/version/downloadVersionJson", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/version/export", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/version/import", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/version/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/version/delete/multi", V2: "POST"},

		// ============ 业务模块 API (888 Admin) ============
		// 容器实例
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/stop", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/start", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/pod/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/log/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/terminal", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/log/stream", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/notebook/auth", V2: "GET"},

		// 训练任务
		{Ptype: "p", V0: "888", V1: "/api/v1/training/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/training/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/training/stop", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/training/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/training/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/training/pod/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/training/log/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/training/log/download", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/training/terminal", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/training/log/stream", V2: "GET"},

		// 推理服务
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/stop", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/start", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/pod/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/log/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/api/key/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/api/key/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/api/key/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/terminal", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/inference/log/stream", V2: "GET"},

		// Image 镜像管理
		{Ptype: "p", V0: "888", V1: "/api/v1/image/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/image/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/image/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/image/delete", V2: "POST"},

		// Product 产品（用户端）
		{Ptype: "p", V0: "888", V1: "/api/v1/product/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/product/get", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/product/filter/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/product/aggregate/list", V2: "GET"},

		// SSH密钥管理
		{Ptype: "p", V0: "888", V1: "/api/v1/sshkey/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/sshkey/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/sshkey/default/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/sshkey/list", V2: "POST"},

		// Secret 管理
		{Ptype: "p", V0: "888", V1: "/api/v1/secret/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/secret/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/secret/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/secret/update", V2: "POST"},

		// PVC 存储卷管理
		{Ptype: "p", V0: "888", V1: "/api/v1/pvc/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/pvc/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/pvc/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/pvc/update", V2: "POST"},

		// Volume 存储管理
		{Ptype: "p", V0: "888", V1: "/api/v1/volume/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/volume/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/volume/expand", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/volume/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/volume/area/list", V2: "GET"},

		// Piper SSH代理管理
		{Ptype: "p", V0: "888", V1: "/api/v1/piper/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/piper/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/piper/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/piper/update", V2: "POST"},

		// Tensorboard 管理
		{Ptype: "p", V0: "888", V1: "/api/v1/tensorboard/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/tensorboard/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/tensorboard/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/tensorboard/update", V2: "POST"},

		// Dashboard
		{Ptype: "p", V0: "888", V1: "/api/v1/dashboard/get", V2: "POST"},

		// Order 财务账单
		{Ptype: "p", V0: "888", V1: "/api/v1/order/overview", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/order/usage/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/order/transaction/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/order/recharge", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/order/order/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/order/invoice/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/order/invoice/apply", V2: "POST"},

		// Account 账号安全
		{Ptype: "p", V0: "888", V1: "/api/v1/account/security/status", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/account/access/log/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/account/active/session/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/account/password/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/account/bind", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/account/mfa/setup", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/account/mfa/activate", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/account/mfa/toggle", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/account/ak/generate", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/account/active/session/kill", V2: "POST"},

		// CMS 产品管理（运营端）
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/product/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/product/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/product/price/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/product/delete", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/product/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/product/detail", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/product/cluster/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/product/area/list", V2: "GET"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/product/node/list", V2: "GET"},
		// 集群管理
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/cluster/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/cluster/add", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/cluster/update", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/cluster/delete", V2: "POST"},

		{Ptype: "p", V0: "888", V1: "/api/v1/cms/node/list", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/node/uncordon", V2: "POST"},
		{Ptype: "p", V0: "888", V1: "/api/v1/cms/node/drain", V2: "POST"},

		// ============ 业务模块 API (8881 User) ============
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/update", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/get", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/stop", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/start", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/pod/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/log/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/terminal", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/log/stream", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/notebook/auth", V2: "GET"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/training/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/training/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/training/stop", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/training/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/training/get", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/training/pod/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/training/log/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/training/log/download", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/training/terminal", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/training/log/stream", V2: "GET"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/stop", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/start", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/list", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/get", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/pod/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/log/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/api/key/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/api/key/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/api/key/list", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/terminal", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/inference/log/stream", V2: "GET"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/image/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/image/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/image/update", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/image/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/product/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/product/get", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/product/filter/list", V2: "GET"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/sshkey/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/sshkey/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/sshkey/default/update", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/sshkey/list", V2: "POST"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/secret/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/secret/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/secret/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/secret/update", V2: "POST"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/pvc/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/pvc/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/pvc/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/pvc/update", V2: "POST"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/volume/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/volume/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/volume/expand", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/volume/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/volume/area/list", V2: "GET"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/piper/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/piper/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/piper/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/piper/update", V2: "POST"},

		// Tensorboard
		{Ptype: "p", V0: "8881", V1: "/api/v1/tensorboard/add", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/tensorboard/delete", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/tensorboard/list", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/tensorboard/update", V2: "POST"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/dashboard/get", V2: "POST"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/order/overview", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/order/usage/list", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/order/transaction/list", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/order/recharge", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/order/order/list", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/order/invoice/list", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/order/invoice/apply", V2: "POST"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/account/security/status", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/account/access/log/list", V2: "POST"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/account/active/session/list", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/account/password/update", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/account/bind", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/account/mfa/setup", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/account/mfa/activate", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/account/mfa/toggle", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/account/ak/generate", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/account/active/session/kill", V2: "POST"},

		{Ptype: "p", V0: "8881", V1: "/api/v1/user/info", V2: "GET"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/user/self/info/update", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/user/self/setting/update", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/user/password/update", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/menu/get", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/menu/tree", V2: "POST"},
		{Ptype: "p", V0: "8881", V1: "/api/v1/jwt/blacklist", V2: "POST"},

		// ============ 初始角色 (9528 Test) ============
		{Ptype: "p", V0: "9528", V1: "/api/v1/user/register", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/api/add", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/api/list", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/api/get", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/api/delete", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/api/update", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/api/all", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/authority/add", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/authority/delete", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/authority/list", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/authority/data/authority/update", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/menu/get", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/menu/list", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/menu/add", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/menu/tree", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/menu/authority/update", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/menu/authority/get", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/menu/delete", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/menu/update", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/menu/get/id", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/user/password/update", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/user/list", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/user/authority/update", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/file/upload", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/file/list", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/file/delete", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/file/edit", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/file/import", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/casbin/update", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/casbin/policy/list", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/jwt/blacklist", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/system/config/get", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/system/config/update", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/customer/add", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/customer/get", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/customer/update", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/customer/delete", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/customer/list", V2: "GET"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/autocode/add", V2: "POST"},
		{Ptype: "p", V0: "9528", V1: "/api/v1/user/info", V2: "GET"},
	}
	// 批量查询已存在的 Casbin Rules，避免循环执行 SQL
	var existingRules []adapter.CasbinRule
	if err := global.GVA_DB.Select("ptype", "v0", "v1", "v2").Find(&existingRules).Error; err != nil {
		return ctx, errors.Wrap(err, "查询已存在Casbin Rules失败")
	}

	existingMap := make(map[string]bool)
	for _, r := range existingRules {
		key := r.Ptype + ":" + r.V0 + ":" + r.V1 + ":" + r.V2
		existingMap[key] = true
	}

	var missingEntities []adapter.CasbinRule
	for _, entity := range entities {
		key := entity.Ptype + ":" + entity.V0 + ":" + entity.V1 + ":" + entity.V2
		if !existingMap[key] {
			missingEntities = append(missingEntities, entity)
		}
	}

	if len(missingEntities) > 0 {
		if err := global.GVA_DB.Create(&missingEntities).Error; err != nil {
			return ctx, errors.Wrap(err, "Casbin 表 ("+i.InitializerName()+") 数据批量初始化失败!")
		}
	}
	e := utils.GetCasbin()
	if err := e.LoadPolicy(); err != nil {
		return nil, errors.Wrapf(err, "CasbinService.UpdateCasbinApi LoadPolicy failed")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
