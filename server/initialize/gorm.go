package initialize

import (
	"gin-vue-admin/global"
	accountModel "gin-vue-admin/model/account"
	"gin-vue-admin/model/cluster"
	"gin-vue-admin/model/dashboard"
	"gin-vue-admin/model/example"
	"gin-vue-admin/model/image"
	"gin-vue-admin/model/inference"
	"gin-vue-admin/model/notebook"
	"gin-vue-admin/model/order"
	"gin-vue-admin/model/podgroup"
	"gin-vue-admin/model/product"
	"gin-vue-admin/model/pvc"
	"gin-vue-admin/model/secret"
	"gin-vue-admin/model/system"
	"gin-vue-admin/model/tensorboard"
	"gin-vue-admin/model/training"
	"os"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
		return GormMysql()
	case "pgsql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Pgsql.Dbname
		return GormPgSql()
	case "oracle":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Oracle.Dbname
		return GormOracle()
	case "mssql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mssql.Dbname
		return GormMssql()
	case "sqlite":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Sqlite.Dbname
		return GormSqlite()
	default:
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
		return GormMysql()
	}
}

func RegisterTables() {
	if global.GVA_CONFIG.System.DisableAutoMigrate {
		logx.Info("auto-migrate is disabled, skipping table registration")
		return
	}

	db := global.GVA_DB
	err := db.AutoMigrate(
		cluster.K8sCluster{},
		notebook.Notebook{},
		notebook.NotebookVolume{},

		pvc.Volume{},
		tensorboard.Tensorboard{},
		secret.SSHKey{},

		order.Order{},
		order.OrderSummary{},
		order.Transaction{},
		order.Wallet{},
		order.Invoice{},
		order.InvoiceTitle{},
		order.InvoiceAddress{},

		image.Image{},
		product.Product{},
		product.ResourceAllocation{},
		podgroup.PodGroup{},

		training.TrainingJob{},
		training.TrainingJobMount{},
		training.TrainingJobEnv{},

		inference.Inference{},
		inference.InferenceMount{},
		inference.InferenceEnv{},
		inference.InferenceApiKeyPolicy{},
		inference.InferenceApiKey{},

		dashboard.DailyMetric{},

		accountModel.AccountAccessLog{},
		accountModel.AccountSecurity{},

		system.SysApi{},
		system.SysIgnoreApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysAutoCodePackage{},
		system.SysExportTemplate{},
		system.Condition{},
		system.JoinTemplate{},
		system.SysParams{},
		system.SysVersion{},
		system.SysError{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},
		example.ExaAttachmentCategory{},
		adapter.CasbinRule{},
	)
	if err != nil {
		logx.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	err = bizModel()

	if err != nil {
		logx.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	logx.Info("register table success")
}
