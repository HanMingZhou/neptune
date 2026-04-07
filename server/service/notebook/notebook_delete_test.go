package notebook

import (
	"context"
	"fmt"
	"testing"

	"gin-vue-admin/global"
	"gin-vue-admin/model/consts"
	nbModel "gin-vue-admin/model/notebook"
	tensorboardModel "gin-vue-admin/model/tensorboard"

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

func TestDeleteNotebookTensorboardResourcesKeepsRecordOnStop(t *testing.T) {
	db := newTestNotebookDB(t)
	svc := &NotebookService{}

	record := &nbModel.Notebook{
		DisplayName:       "demo",
		InstanceName:      "notebook-abc123",
		Namespace:         "zzz",
		Status:            "STOPPED",
		EnableTensorboard: true,
	}
	require.NoError(t, db.Create(record).Error)

	tbRecord := &tensorboardModel.Tensorboard{
		InstanceName: "notebook-abc123-tb",
		OwnerType:    consts.NotebookInstance,
		OwnerID:      record.ID,
		Namespace:    record.Namespace,
		LogsPath:     "pvc://notebook-abc123-workspace/logs",
		Status:       "RUNNING",
	}
	require.NoError(t, db.Create(tbRecord).Error)

	require.NoError(t, svc.deleteNotebookTensorboardResources(context.Background(), record, nil, false))

	var countAfterStop int64
	require.NoError(t, db.Model(&tensorboardModel.Tensorboard{}).Where("owner_id = ? AND owner_type = ?", record.ID, consts.NotebookInstance).Count(&countAfterStop).Error)
	require.Equal(t, int64(1), countAfterStop)

	require.NoError(t, svc.deleteNotebookTensorboardResources(context.Background(), record, nil, true))

	var countAfterDelete int64
	require.NoError(t, db.Model(&tensorboardModel.Tensorboard{}).Where("owner_id = ? AND owner_type = ?", record.ID, consts.NotebookInstance).Count(&countAfterDelete).Error)
	require.Equal(t, int64(0), countAfterDelete)
}

func TestUpsertNotebookTensorboardRecordReusesExistingRow(t *testing.T) {
	db := newTestNotebookDB(t)
	svc := &NotebookService{}

	record := &nbModel.Notebook{
		DisplayName:       "demo",
		InstanceName:      "notebook-abc123",
		Namespace:         "zzz",
		Status:            "STOPPED",
		UserId:            7,
		ClusterID:         9,
		EnableTensorboard: true,
	}
	require.NoError(t, db.Create(record).Error)

	tbRecord := &tensorboardModel.Tensorboard{
		InstanceName: "legacy-notebook-tb",
		OwnerType:    consts.NotebookInstance,
		OwnerID:      record.ID,
		Namespace:    record.Namespace,
		LogsPath:     "pvc://legacy/logs",
		Status:       "STOPPED",
		UserId:       1,
		ClusterID:    2,
	}
	require.NoError(t, db.Create(tbRecord).Error)

	record.TensorboardID = tbRecord.ID

	updated, err := svc.upsertNotebookTensorboardRecord(record, "notebook-abc123-tb", "pvc://notebook-abc123-workspace/logs")
	require.NoError(t, err)
	require.Equal(t, tbRecord.ID, updated.ID)

	var count int64
	require.NoError(t, db.Model(&tensorboardModel.Tensorboard{}).Where("owner_id = ? AND owner_type = ?", record.ID, consts.NotebookInstance).Count(&count).Error)
	require.Equal(t, int64(1), count)

	var refreshed tensorboardModel.Tensorboard
	require.NoError(t, db.First(&refreshed, tbRecord.ID).Error)
	require.Equal(t, "notebook-abc123-tb", refreshed.InstanceName)
	require.Equal(t, "pvc://notebook-abc123-workspace/logs", refreshed.LogsPath)
	require.Equal(t, record.UserId, refreshed.UserId)
	require.Equal(t, record.ClusterID, refreshed.ClusterID)
}

func newTestNotebookDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&nbModel.Notebook{}, &nbModel.NotebookVolume{}, &tensorboardModel.Tensorboard{}))

	previous := global.GVA_DB
	global.GVA_DB = db
	t.Cleanup(func() {
		global.GVA_DB = previous
	})

	return db
}
