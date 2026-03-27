package podgroup

import (
	"context"
	"fmt"
	"gin-vue-admin/model/consts"

	"github.com/zeromicro/go-zero/core/logx"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vcv1beta1 "volcano.sh/apis/pkg/apis/scheduling/v1beta1"
	vcclientset "volcano.sh/apis/pkg/client/clientset/versioned"
)

// PodGroupManager PodGroup 管理器
// 注意：PodGroup 由 Volcano 调度器自动创建，本管理器主要用于查询和删除操作
type PodGroupManager struct {
	vcClient vcclientset.Interface
}

// NewPodGroupManager 创建 PodGroup 管理器
func NewPodGroupManager(vcClient vcclientset.Interface) *PodGroupManager {
	return &PodGroupManager{
		vcClient: vcClient,
	}
}

// GetPodGroup 获取 PodGroup
func (m *PodGroupManager) GetPodGroup(ctx context.Context, name, namespace string) (*vcv1beta1.PodGroup, error) {
	return m.vcClient.SchedulingV1beta1().PodGroups(namespace).Get(
		ctx,
		name,
		metav1.GetOptions{},
	)
}

// DeletePodGroup 删除 PodGroup
// 通常在删除 Notebook 后，Volcano 会自动清理 PodGroup，但也可以手动删除
func (m *PodGroupManager) DeletePodGroup(ctx context.Context, name, namespace string) error {
	err := m.vcClient.SchedulingV1beta1().PodGroups(namespace).Delete(
		ctx,
		name,
		metav1.DeleteOptions{},
	)
	if err != nil {
		logx.Error("删除 PodGroup 失败",
			logx.Field("name", name),
			logx.Field("namespace", namespace),
			logx.Field("error", err))
		return err
	}

	logx.Info("删除 PodGroup 成功",
		logx.Field("name", name),
		logx.Field("namespace", namespace))

	return nil
}

// ListPodGroups 列出 PodGroup
func (m *PodGroupManager) ListPodGroups(ctx context.Context, namespace string) (*vcv1beta1.PodGroupList, error) {
	return m.vcClient.SchedulingV1beta1().PodGroups(namespace).List(
		ctx,
		metav1.ListOptions{},
	)
}

// VolcanoSchedulerName Volcano 调度器名称
const VolcanoSchedulerName = "volcano"

// VolcanoPodGroupAnnotation PodGroup 名称的 annotation key
// Volcano 调度器会根据此 annotation 自动创建 PodGroup
const VolcanoPodGroupAnnotation = "scheduling.volcano.sh/group-name"

// VolcanoQueueAnnotation 队列名称的 annotation key
const VolcanoQueueAnnotation = "volcano.sh/queue-name"

// BuildVolcanoAnnotations 构建 Volcano 调度器需要的 annotations
// instanceName: 实例名称，用作 PodGroup 名称
// queue: 队列名称（可选，默认为 "default"）
func BuildVolcanoAnnotations(instanceName string, queue string) map[string]string {
	annotations := map[string]string{
		VolcanoPodGroupAnnotation: instanceName,
	}
	if queue != "" {
		annotations[VolcanoQueueAnnotation] = queue
	}
	return annotations
}

// BuildVolcanoLabels 构建labels，用于 Informer 识别
func BuildVolcanoLabels(instanceName, instanceType, ownerType string, jobID uint) map[string]string {
	labels := map[string]string{
		consts.LabelApp:          instanceName,
		consts.LabelInstanceType: instanceType,
		consts.LabelOwnerType:    ownerType,
	}
	if jobID > 0 {
		labels[consts.LabelJobID] = fmt.Sprintf("%d", jobID)
	}
	return labels
}
