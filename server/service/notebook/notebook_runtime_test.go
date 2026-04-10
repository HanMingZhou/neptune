package notebook

import (
	"testing"

	nbModel "gin-vue-admin/model/notebook"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestBuildNotebookNormalizesDuplicateVolumeMounts(t *testing.T) {
	nbRef := &nbModel.Notebook{
		InstanceName: "notebook-a42a9b",
		Namespace:    "zzz",
		Image:        "dockerhub.kubekey.local/notebook/jupyter:test",
		CPU:          12,
		Memory:       50,
		GPU:          1,
		StorageSize:  20,
		VolumeMounts: []nbModel.NotebookVolume{
			{
				Name:       nbModel.Workspace,
				MountsPath: "/home/notebook",
				Type:       nbModel.Workspace,
				PVCName:    "notebook-97c140-workspace",
				Size:       10,
			},
			{
				Name:       "hmz",
				MountsPath: "/home/notebook/123",
				Type:       nbModel.VolumeTypeDataset,
				PVCId:      1,
				PVCName:    "vol-1767888821-dqcm",
			},
			{
				Name:       nbModel.Workspace,
				MountsPath: nbModel.DefaultWorkspacePath,
				Type:       nbModel.Workspace,
				PVCName:    "notebook-a42a9b-workspace",
				Size:       20,
			},
			{
				Name:       "hmz",
				MountsPath: "/home/notebook/123",
				Type:       nbModel.VolumeTypeDataset,
				PVCId:      1,
				PVCName:    "vol-1767888821-dqcm",
			},
			{
				Name:       "111",
				MountsPath: "/home/notebook/neptune",
				Type:       nbModel.VolumeTypeDataset,
				PVCId:      2,
				PVCName:    "vol-1775182658-plxm",
			},
		},
	}

	notebookObj := buildNotebook(nbRef, "")
	container := notebookObj.Spec.Template.Spec.Containers[0]

	volumeClaims := make(map[string]string)
	workspaceVolume := corev1.Volume{}
	for _, volume := range notebookObj.Spec.Template.Spec.Volumes {
		if volume.Name == nbModel.Workspace {
			workspaceVolume = volume
		}
		if volume.PersistentVolumeClaim == nil {
			continue
		}
		volumeClaims[volume.Name] = volume.PersistentVolumeClaim.ClaimName
	}

	mountPaths := make(map[string]string)
	for _, mount := range container.VolumeMounts {
		mountPaths[mount.Name] = mount.MountPath
	}

	require.NotNil(t, workspaceVolume.EmptyDir)
	require.NotNil(t, workspaceVolume.EmptyDir.SizeLimit)
	require.Equal(t, "20Gi", workspaceVolume.EmptyDir.SizeLimit.String())
	require.Equal(t, nbModel.DefaultWorkspacePath, mountPaths[nbModel.Workspace])
	require.Equal(t, "vol-1767888821-dqcm", volumeClaims["hmz"])
	require.Equal(t, "/home/notebook/123", mountPaths["hmz"])
	require.Equal(t, "vol-1775182658-plxm", volumeClaims["111"])
	require.Equal(t, "/home/notebook/neptune", mountPaths["111"])
	require.Len(t, volumeClaims, 2)
	require.Len(t, mountPaths, 4)
}

func TestGetNotebookTensorboardLogsPathRequiresPersistentMount(t *testing.T) {
	svc := &NotebookService{}
	nbRef := &nbModel.Notebook{
		VolumeMounts: []nbModel.NotebookVolume{
			{
				Name:       nbModel.Workspace,
				MountsPath: nbModel.DefaultWorkspacePath,
				Type:       nbModel.Workspace,
				Size:       50,
			},
			{
				Name:       "dataset",
				MountsPath: "/data",
				Type:       nbModel.VolumeTypeDataset,
				PVCId:      1,
				PVCName:    "vol-1",
			},
		},
	}

	nbRef.TensorboardLogPath = "logs"
	require.Equal(t, "", svc.getNotebookTensorboardLogsPath(nbRef))

	nbRef.TensorboardLogPath = "/data/tensorboard/run1"
	require.Equal(t, "pvc://vol-1/tensorboard/run1", svc.getNotebookTensorboardLogsPath(nbRef))
}
