package notebook

import (
	"testing"

	nbModel "gin-vue-admin/model/notebook"

	"github.com/stretchr/testify/require"
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
	for _, volume := range notebookObj.Spec.Template.Spec.Volumes {
		if volume.PersistentVolumeClaim == nil {
			continue
		}
		volumeClaims[volume.Name] = volume.PersistentVolumeClaim.ClaimName
	}

	mountPaths := make(map[string]string)
	for _, mount := range container.VolumeMounts {
		mountPaths[mount.Name] = mount.MountPath
	}

	require.Equal(t, "notebook-a42a9b-workspace", volumeClaims[nbModel.Workspace])
	require.Equal(t, nbModel.DefaultWorkspacePath, mountPaths[nbModel.Workspace])
	require.Equal(t, "vol-1767888821-dqcm", volumeClaims["hmz"])
	require.Equal(t, "/home/notebook/123", mountPaths["hmz"])
	require.Equal(t, "vol-1775182658-plxm", volumeClaims["111"])
	require.Equal(t, "/home/notebook/neptune", mountPaths["111"])
	require.Len(t, volumeClaims, 3)
	require.Len(t, mountPaths, 4)
}
