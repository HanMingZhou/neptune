package initialize

import (
	"context"
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
	"path/filepath"
	"sort"
	"strings"

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
		product.ProductPrice{},
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

	if err = product.SyncLegacyComputePriceItems(context.Background(), db); err != nil {
		logx.Error("sync legacy product prices failed", zap.Error(err))
		os.Exit(0)
	}

	err = bizModel()

	if err != nil {
		logx.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	logx.Info("register table success")
}

func LoadInitSQLFiles() error {
	sqlDir, err := findSQLDir()
	if err != nil {
		logx.Error("find sql dir failed", zap.Error(err))
		os.Exit(0)
	}

	entries, err := os.ReadDir(sqlDir)
	if err != nil {
		logx.Error("read sql dir failed", zap.Error(err))
		os.Exit(0)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".sql") {
			continue
		}

		filePath := filepath.Join(sqlDir, entry.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			logx.Error("read sql file failed", zap.Error(err))
			os.Exit(0)
		}

		statements := splitSQLStatements(string(content))
		for _, statement := range statements {
			if err := global.GVA_DB.Exec(statement).Error; err != nil {
				logx.Error("exec sql file failed", zap.Error(err))
				os.Exit(0)
			}
		}

		global.GVA_LOG.Info("初始化 SQL 文件加载完成", zap.String("file", filePath), zap.Int("statements", len(statements)))
	}

	return nil
}

func findSQLDir() (string, error) {
	candidates := []string{
		filepath.Join("deploy", "sql"),
		filepath.Join("..", "deploy", "sql"),
	}

	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(wd, "deploy", "sql"),
			filepath.Join(wd, "..", "deploy", "sql"),
		)
	}

	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		candidates = append(candidates,
			filepath.Join(execDir, "deploy", "sql"),
			filepath.Join(execDir, "..", "deploy", "sql"),
		)
	}

	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
	}
	return "", nil
}

func splitSQLStatements(content string) []string {
	var (
		builder         strings.Builder
		statements      []string
		inSingleQuote   bool
		inDoubleQuote   bool
		inBacktickQuote bool
		inLineComment   bool
		inBlockComment  bool
		previous        rune
	)

	for _, current := range content {
		nextTwo := string([]rune{previous, current})

		if inLineComment {
			builder.WriteRune(current)
			if current == '\n' {
				inLineComment = false
			}
			previous = current
			continue
		}

		if inBlockComment {
			builder.WriteRune(current)
			if nextTwo == "*/" {
				inBlockComment = false
			}
			previous = current
			continue
		}

		if !inSingleQuote && !inDoubleQuote && !inBacktickQuote {
			switch {
			case previous == '-' && current == '-':
				inLineComment = true
			case previous == '/' && current == '*':
				inBlockComment = true
			}
		}

		switch current {
		case '\'':
			if !inDoubleQuote && !inBacktickQuote && previous != '\\' {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote && !inBacktickQuote && previous != '\\' {
				inDoubleQuote = !inDoubleQuote
			}
		case '`':
			if !inSingleQuote && !inDoubleQuote {
				inBacktickQuote = !inBacktickQuote
			}
		case ';':
			if !inSingleQuote && !inDoubleQuote && !inBacktickQuote {
				statement := strings.TrimSpace(builder.String())
				if statement != "" {
					statements = append(statements, statement)
				}
				builder.Reset()
				previous = 0
				continue
			}
		}

		builder.WriteRune(current)
		previous = current
	}

	if statement := strings.TrimSpace(builder.String()); statement != "" {
		statements = append(statements, statement)
	}

	return statements
}
