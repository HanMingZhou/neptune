package pvc

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	clusterModel "gin-vue-admin/model/cluster"
	pvcReq "gin-vue-admin/model/pvc/request"

	"github.com/zeromicro/go-zero/core/logx"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PVCManager interface {
	CreatePVCs(ctx context.Context, pvc *pvcReq.AddPVCReq, instanceType string) error
	DeletePVCs(ctx context.Context, pvc *pvcReq.DeletePVCReq, instanceType string) error
}

var _ PVCManager = (*K8sPVCService)(nil)

type K8sPVCService struct {
	client kubernetes.Interface
}

func NewK8sPVCManager(client kubernetes.Interface) *K8sPVCService {
	return &K8sPVCService{client: client}
}

func (m *K8sPVCService) CreatePVCs(ctx context.Context, pvc *pvcReq.AddPVCReq, instanceType string) error {
	// 查询集群配置获取默认 StorageClass
	var storageClassName string
	var cluster clusterModel.K8sCluster
	if err := global.GVA_DB.Where("id = ?", pvc.ClusterId).First(&cluster).Error; err == nil {
		if cluster.StorageClass != "" {
			storageClassName = cluster.StorageClass
		}
	}

	for _, mount := range pvc.VolumeMounts {
		if mount.PVCId > 0 {
			continue // 如果是已有的 PVC，跳过创建
		}
		// 构造 PVC 名称: <notebook-name>-<volume-name>
		// 例如: demo-workspace
		pvcObj := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      mount.PVCName,
				Namespace: pvc.Namespace,
				Labels: map[string]string{
					"app":        pvc.InstanceName,
					instanceType: pvc.InstanceName, // 用于级联删除
				},
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse(fmt.Sprintf("%dGi", mount.Size)),
					},
				},
				StorageClassName: &storageClassName,
			},
		}

		if _, err := m.client.CoreV1().
			PersistentVolumeClaims(pvc.Namespace).
			Create(ctx, pvcObj, metav1.CreateOptions{}); err != nil {
			logx.Error("创建PVC失败", err)
			return err
		}
	}
	return nil
}

func (m *K8sPVCService) DeletePVCs(ctx context.Context, pvc *pvcReq.DeletePVCReq, instanceType string) error {
	// 根据 Label 删除所有关联的 PVC
	listOpts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", instanceType, pvc.InstanceName),
	}
	if err := m.client.CoreV1().
		PersistentVolumeClaims(pvc.Namespace).
		DeleteCollection(ctx, metav1.DeleteOptions{}, listOpts); err != nil {
		logx.Error("删除关联PVC失败", err)
		return err
	}
	return nil
}
