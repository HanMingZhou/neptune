package pvc

import (
	"context"
	"fmt"
	"math"

	"gin-vue-admin/global"
	pvcModel "gin-vue-admin/model/pvc"

	"github.com/zeromicro/go-zero/core/logx"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const bytesPerGi = 1024 * 1024 * 1024

// VolumeRuntimeState 表示 PVC 当前在集群中的真实状态。
// ActualSize 取自 status.capacity，RequestedSize 取自 spec.requests.storage。
type VolumeRuntimeState struct {
	ActualSize    int64
	RequestedSize int64
	ResizePending bool
}

// GetVolumeRuntimeState 从集群读取 PVC 当前的实际容量和已申请容量。
func GetVolumeRuntimeState(ctx context.Context, clusterID uint, namespace, pvcName string, fallbackSize int64) (*VolumeRuntimeState, error) {
	state := &VolumeRuntimeState{
		ActualSize:    fallbackSize,
		RequestedSize: fallbackSize,
	}

	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterID)
	if cluster == nil {
		return state, fmt.Errorf("对应集群不可用")
	}

	pvc, err := cluster.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, metav1.GetOptions{})
	if err != nil {
		return state, err
	}

	requestedSize := quantityToGi(pvc.Spec.Resources.Requests[corev1.ResourceStorage])
	actualSize := requestedSize
	if capacity, ok := pvc.Status.Capacity[corev1.ResourceStorage]; ok {
		if size := quantityToGi(capacity); size > 0 {
			actualSize = size
		}
	}

	if actualSize <= 0 {
		actualSize = fallbackSize
	}
	if requestedSize <= 0 {
		requestedSize = actualSize
	}

	state.ActualSize = actualSize
	state.RequestedSize = requestedSize
	state.ResizePending = requestedSize > actualSize
	return state, nil
}

// SyncVolumeRuntimeState 会以集群中的真实容量校准数据库 size 字段。
// 数据库仅保存已经实际生效的容量，避免出现“申请扩容成功但底层尚未扩到位”时的假数据。
func SyncVolumeRuntimeState(ctx context.Context, volumeID, clusterID uint, namespace, pvcName string, currentSize int64) (*VolumeRuntimeState, error) {
	state, err := GetVolumeRuntimeState(ctx, clusterID, namespace, pvcName, currentSize)
	if err != nil {
		return state, err
	}

	if volumeID == 0 || state.ActualSize <= 0 || state.ActualSize == currentSize {
		return state, nil
	}

	if err := global.GVA_DB.Model(&pvcModel.Volume{}).Where("id = ?", volumeID).Update("size", state.ActualSize).Error; err != nil {
		logx.Error("同步Volume容量到数据库失败",
			logx.Field("volumeId", volumeID),
			logx.Field("actualSize", state.ActualSize),
			logx.Field("err", err))
	}

	return state, nil
}

func quantityToGi(quantity resource.Quantity) int64 {
	if quantity.IsZero() {
		return 0
	}

	return int64(math.Ceil(float64(quantity.Value()) / float64(bytesPerGi)))
}
