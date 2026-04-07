package notebook

import (
	"testing"
	"time"

	"gin-vue-admin/model/consts"

	nbv1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNotebookStatusFromResource(t *testing.T) {
	t.Run("deleting notebook stays deleting", func(t *testing.T) {
		now := metav1.NewTime(time.Now())
		notebookObj := &nbv1.Notebook{
			ObjectMeta: metav1.ObjectMeta{
				DeletionTimestamp: &now,
			},
		}

		require.Equal(t, consts.NotebookStatusDeleting, notebookStatusFromResource(notebookObj))
	})

	t.Run("running container is running", func(t *testing.T) {
		notebookObj := &nbv1.Notebook{
			Status: nbv1.NotebookStatus{
				ContainerState: corev1.ContainerState{
					Running: &corev1.ContainerStateRunning{},
				},
			},
		}

		require.Equal(t, consts.NotebookStatusRunning, notebookStatusFromResource(notebookObj))
	})

	t.Run("waiting container is pending", func(t *testing.T) {
		notebookObj := &nbv1.Notebook{
			Status: nbv1.NotebookStatus{
				ContainerState: corev1.ContainerState{
					Waiting: &corev1.ContainerStateWaiting{Reason: "ContainerCreating"},
				},
			},
		}

		require.Equal(t, "Waiting: ContainerCreating", notebookStatusFromResource(notebookObj))
	})

	t.Run("empty status defaults to pending when resource exists", func(t *testing.T) {
		notebookObj := &nbv1.Notebook{}

		require.Equal(t, consts.NotebookStatusPending, notebookStatusFromResource(notebookObj))
	})
}
