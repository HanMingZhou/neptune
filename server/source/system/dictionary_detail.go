package system

import (
	"context"
	"gin-vue-admin/global"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/service/system"

	"github.com/pkg/errors"
)

const initOrderDictDetail = initOrderDict + 1

type initDictDetail struct{}

// auto run
func init() {
	system.RegisterInit(initOrderDictDetail, &initDictDetail{})
}

func (i *initDictDetail) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, global.GVA_DB.AutoMigrate(&sysModel.SysDictionaryDetail{})
}

func (i *initDictDetail) TableCreated(ctx context.Context) bool {
	return global.GVA_DB.Migrator().HasTable(&sysModel.SysDictionaryDetail{})
}

func (i *initDictDetail) InitializerName() string {
	return sysModel.SysDictionaryDetail{}.TableName()
}

func (i *initDictDetail) InitializeData(ctx context.Context) (context.Context, error) {
	True := true
	var dicts []sysModel.SysDictionary

	// 1. 尝试从 Context 获取 (由 initDict 填充)
	if val := ctx.Value(new(initDict).InitializerName()); val != nil {
		if ds, ok := val.([]sysModel.SysDictionary); ok {
			dicts = ds
		}
	}

	// 2. 如果 Context 中没有或者不全（防御性编程），则查库补全
	// 顺序必须与下方 hardcoded 的索引一致:
	// 0: gender, 1: int, 2: time.Time, 3: float64, 4: string, 5: bool
	if len(dicts) < 6 {
		types := []string{"gender", "int", "time.Time", "float64", "string", "bool"}
		var dbDicts []sysModel.SysDictionary
		if err := global.GVA_DB.Where("type IN ?", types).Find(&dbDicts).Error; err != nil {
			return ctx, errors.Wrap(err, "查询字典失败")
		}

		dictMap := make(map[string]sysModel.SysDictionary)
		for _, d := range dbDicts {
			dictMap[d.Type] = d
		}

		dicts = make([]sysModel.SysDictionary, 6)
		for i, t := range types {
			if d, ok := dictMap[t]; ok {
				dicts[i] = d
			}
		}
	}

	if len(dicts) > 0 {
		dicts[0].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
			{Label: "男", Value: "1", Status: &True, Sort: 1},
			{Label: "女", Value: "2", Status: &True, Sort: 2},
		}
	}

	if len(dicts) > 1 {
		dicts[1].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
			{Label: "smallint", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
			{Label: "mediumint", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
			{Label: "int", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
			{Label: "bigint", Value: "4", Status: &True, Extend: "mysql", Sort: 4},
			{Label: "int2", Value: "5", Status: &True, Extend: "pgsql", Sort: 5},
			{Label: "int4", Value: "6", Status: &True, Extend: "pgsql", Sort: 6},
			{Label: "int6", Value: "7", Status: &True, Extend: "pgsql", Sort: 7},
			{Label: "int8", Value: "8", Status: &True, Extend: "pgsql", Sort: 8},
		}
	}

	if len(dicts) > 2 {
		dicts[2].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
			{Label: "date", Value: "0", Status: &True, Extend: "mysql", Sort: 0},
			{Label: "time", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
			{Label: "year", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
			{Label: "datetime", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
			{Label: "timestamp", Value: "5", Status: &True, Extend: "mysql", Sort: 5},
			{Label: "timestamptz", Value: "6", Status: &True, Extend: "pgsql", Sort: 5},
		}
	}

	if len(dicts) > 3 {
		dicts[3].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
			{Label: "float", Value: "0", Status: &True, Extend: "mysql", Sort: 0},
			{Label: "double", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
			{Label: "decimal", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
			{Label: "numeric", Value: "3", Status: &True, Extend: "pgsql", Sort: 3},
			{Label: "smallserial", Value: "4", Status: &True, Extend: "pgsql", Sort: 4},
		}
	}

	if len(dicts) > 4 {
		dicts[4].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
			{Label: "char", Value: "0", Status: &True, Extend: "mysql", Sort: 0},
			{Label: "varchar", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
			{Label: "tinyblob", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
			{Label: "tinytext", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
			{Label: "text", Value: "4", Status: &True, Extend: "mysql", Sort: 4},
			{Label: "blob", Value: "5", Status: &True, Extend: "mysql", Sort: 5},
			{Label: "mediumblob", Value: "6", Status: &True, Extend: "mysql", Sort: 6},
			{Label: "mediumtext", Value: "7", Status: &True, Extend: "mysql", Sort: 7},
			{Label: "longblob", Value: "8", Status: &True, Extend: "mysql", Sort: 8},
			{Label: "longtext", Value: "9", Status: &True, Extend: "mysql", Sort: 9},
		}
	}

	if len(dicts) > 5 {
		dicts[5].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
			{Label: "tinyint", Value: "1", Extend: "mysql", Status: &True},
			{Label: "bool", Value: "2", Extend: "pgsql", Status: &True},
		}
	}

	for _, dict := range dicts {
		// dict.ID 必须存在，否则更新关联失败
		if dict.ID != 0 {
			if err := global.GVA_DB.Model(&dict).Association("SysDictionaryDetails").
				Replace(dict.SysDictionaryDetails); err != nil {
				return ctx, errors.Wrap(err, sysModel.SysDictionaryDetail{}.TableName()+"表数据初始化失败!")
			}
		}
	}
	return ctx, nil
}
