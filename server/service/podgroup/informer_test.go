package podgroup

import (
	"testing"

	"gin-vue-admin/model/consts"
	podgroupModel "gin-vue-admin/model/podgroup"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vcv1beta1 "volcano.sh/apis/pkg/apis/scheduling/v1beta1"
)

func TestIdentifyResourceUsesExplicitLabels(t *testing.T) {
	factory := &PodGroupInformerFactory{}
	pg := &vcv1beta1.PodGroup{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				consts.LabelInstanceType: consts.TrainingInstance,
				consts.LabelJobID:        "42",
				consts.LabelApp:          "training-abc123",
			},
		},
	}

	res := factory.identifyResource(pg)
	require.True(t, res.Found())
	require.Equal(t, consts.TrainingInstance, res.InstanceType)
	require.Equal(t, uint(42), res.OwnerID)
	require.Equal(t, "training-abc123", res.InstanceName)
}

func TestResolveDeleteResourceFallsBackToStoredPodGroup(t *testing.T) {
	db := newTestPodGroupDB(t)
	require.NoError(t, db.Create(&podgroupModel.PodGroup{
		Name:         "training-a1b2c3",
		Namespace:    "default",
		InstanceName: "training-a1b2c3",
		InstanceType: consts.TrainingInstance,
		OwnerID:      88,
		OwnerType:    consts.TrainingInstance,
		UserID:       1,
		ClusterID:    1,
	}).Error)

	factory := &PodGroupInformerFactory{db: db}
	pg := &vcv1beta1.PodGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "training-a1b2c3",
			Namespace: "default",
		},
	}

	res := factory.resolveDeleteResource(pg)
	require.True(t, res.Found())
	require.Equal(t, consts.TrainingInstance, res.InstanceType)
	require.Equal(t, uint(88), res.OwnerID)
	require.Equal(t, "training-a1b2c3", res.InstanceName)
}

func newTestPodGroupDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&podgroupModel.PodGroup{}))
	return db
}
