package global

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	vcv1beta1 "volcano.sh/apis/pkg/apis/scheduling/v1beta1"
	vcclientset "volcano.sh/apis/pkg/client/clientset/versioned"
)

// VolcanoClient 封装 Volcano CRD 操作
type VolcanoClient struct {
	client vcclientset.Interface
}

// NewVolcanoClient 创建 Volcano 客户端
func NewVolcanoClient(config *rest.Config) (*VolcanoClient, error) {
	client, err := vcclientset.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建 volcano client 失败: %w", err)
	}
	return &VolcanoClient{client: client}, nil
}

// Clientset 获取底层的 Clientset (用于 Informer 等)
func (c *VolcanoClient) Clientset() vcclientset.Interface {
	return c.client
}

// ======== PodGroup ========
// CreatePodGroup 创建 PodGroup
func (c *VolcanoClient) CreatePodGroup(ctx context.Context, pg *vcv1beta1.PodGroup) (*vcv1beta1.PodGroup, error) {
	return c.client.SchedulingV1beta1().PodGroups(pg.Namespace).Create(ctx, pg, metav1.CreateOptions{})
}

// GetPodGroup 获取 PodGroup
func (c *VolcanoClient) GetPodGroup(ctx context.Context, namespace, name string) (*vcv1beta1.PodGroup, error) {
	return c.client.SchedulingV1beta1().PodGroups(namespace).Get(ctx, name, metav1.GetOptions{})
}

// DeletePodGroup 删除 PodGroup
func (c *VolcanoClient) DeletePodGroup(ctx context.Context, namespace, name string) error {
	return c.client.SchedulingV1beta1().PodGroups(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// ListPodGroups 列出 PodGroup
func (c *VolcanoClient) ListPodGroups(ctx context.Context, namespace string) (*vcv1beta1.PodGroupList, error) {
	return c.client.SchedulingV1beta1().PodGroups(namespace).List(ctx, metav1.ListOptions{})
}

// ======== Queue ========
// CreateQueue 创建 Queue
func (c *VolcanoClient) CreateQueue(ctx context.Context, queue *vcv1beta1.Queue) (*vcv1beta1.Queue, error) {
	return c.client.SchedulingV1beta1().Queues().Create(ctx, queue, metav1.CreateOptions{})
}

// GetQueue 获取 Queue
func (c *VolcanoClient) GetQueue(ctx context.Context, name string) (*vcv1beta1.Queue, error) {
	return c.client.SchedulingV1beta1().Queues().Get(ctx, name, metav1.GetOptions{})
}

// DeleteQueue 删除 Queue
func (c *VolcanoClient) DeleteQueue(ctx context.Context, name string) error {
	return c.client.SchedulingV1beta1().Queues().Delete(ctx, name, metav1.DeleteOptions{})
}

// ListQueues 列出 Queue
func (c *VolcanoClient) ListQueues(ctx context.Context) (*vcv1beta1.QueueList, error) {
	return c.client.SchedulingV1beta1().Queues().List(ctx, metav1.ListOptions{})
}
