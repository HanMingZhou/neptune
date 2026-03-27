package system

import (
	"context"
	"gin-vue-admin/global"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/service/system"

	"github.com/pkg/errors"
)

const initOrderDict = initOrderCasbin + 1

type initDict struct{}

// auto run
func init() {
	system.RegisterInit(initOrderDict, &initDict{})
}

func (i *initDict) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, global.GVA_DB.AutoMigrate(&sysModel.SysDictionary{})
}

func (i *initDict) TableCreated(ctx context.Context) bool {
	return global.GVA_DB.Migrator().HasTable(&sysModel.SysDictionary{})
}

func (i *initDict) InitializerName() string {
	return sysModel.SysDictionary{}.TableName()
}

func (i *initDict) InitializeData(ctx context.Context) (next context.Context, err error) {
	True := true
	entities := []sysModel.SysDictionary{
		{Name: "性别", Type: "gender", Status: &True, Desc: "性别字典"},
		{Name: "数据库int类型", Type: "int", Status: &True, Desc: "int类型对应的数据库类型"},
		{Name: "数据库时间日期类型", Type: "time.Time", Status: &True, Desc: "数据库时间日期类型"},
		{Name: "数据库浮点型", Type: "float64", Status: &True, Desc: "数据库浮点型"},
		{Name: "数据库字符串", Type: "string", Status: &True, Desc: "数据库字符串"},
		{Name: "数据库bool类型", Type: "bool", Status: &True, Desc: "数据库bool类型"},
	}

	// 查询已存在的 Dictionary
	var existingDicts []sysModel.SysDictionary
	if err := global.GVA_DB.Find(&existingDicts).Error; err != nil {
		return ctx, errors.Wrap(err, "查询已存在Dictionary失败")
	}

	existingMap := make(map[string]sysModel.SysDictionary)
	for _, d := range existingDicts {
		existingMap[d.Type] = d
	}

	var missingEntities []sysModel.SysDictionary
	for _, entity := range entities {
		if _, ok := existingMap[entity.Type]; !ok {
			missingEntities = append(missingEntities, entity)
		}
	}

	if len(missingEntities) > 0 {
		if err = global.GVA_DB.Create(&missingEntities).Error; err != nil {
			return ctx, errors.Wrap(err, sysModel.SysDictionary{}.TableName()+"表数据批量初始化失败!")
		}
		// 将新创建的补充到 map
		for _, d := range missingEntities {
			existingMap[d.Type] = d
		}
	}

	// 必须回填 ID 到 entities，因为后续步骤 (sys_dictionary_detail) 依赖 Context 中的 ID 来建立关联
	for i := range entities {
		if d, ok := existingMap[entities[i].Type]; ok {
			entities[i].ID = d.ID
		}
	}

	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}
