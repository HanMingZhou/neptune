package notebook

import (
	"fmt"
	"testing"

	"gin-vue-admin/global"
	nbModel "gin-vue-admin/model/notebook"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSoftDeleteNotebookRecordAlsoDeletesVolumes(t *testing.T) {
	db := newTestNotebookDB(t)
	svc := &NotebookService{}

	record := &nbModel.Notebook{
		DisplayName:  "demo",
		InstanceName: "notebook-abc123",
		Namespace:    "zzz",
		Status:       "RUNNING",
	}
	require.NoError(t, db.Create(record).Error)
	require.NoError(t, db.Create([]nbModel.NotebookVolume{
		{
			NotebookID: record.ID,
			Name:       "workspace",
			MountsPath: "/home/system",
			Type:       nbModel.Workspace,
		},
		{
			NotebookID: record.ID,
			Name:       "dataset",
			MountsPath: "/data",
			Type:       nbModel.VolumeTypeDataset,
		},
	}).Error)

	require.NoError(t, svc.softDeleteNotebookRecord(record.ID))

	var notebookCount int64
	require.NoError(t, db.Model(&nbModel.Notebook{}).Where("id = ?", record.ID).Count(&notebookCount).Error)
	require.Equal(t, int64(0), notebookCount)

	var volumeCount int64
	require.NoError(t, db.Model(&nbModel.NotebookVolume{}).Where("notebook_id = ?", record.ID).Count(&volumeCount).Error)
	require.Equal(t, int64(0), volumeCount)

	var deletedNotebookCount int64
	require.NoError(t, db.Unscoped().Model(&nbModel.Notebook{}).Where("id = ?", record.ID).Count(&deletedNotebookCount).Error)
	require.Equal(t, int64(1), deletedNotebookCount)

	var deletedVolumeCount int64
	require.NoError(t, db.Unscoped().Model(&nbModel.NotebookVolume{}).Where("notebook_id = ?", record.ID).Count(&deletedVolumeCount).Error)
	require.Equal(t, int64(2), deletedVolumeCount)
}

func newTestNotebookDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&nbModel.Notebook{}, &nbModel.NotebookVolume{}))

	previous := global.GVA_DB
	global.GVA_DB = db
	t.Cleanup(func() {
		global.GVA_DB = previous
	})

	return db
}
