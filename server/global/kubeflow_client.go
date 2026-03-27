package global

import (
	"context"

	nbv1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1"
	tbv1alpha1 "github.com/kubeflow/kubeflow/components/tensorboard-controller/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// NotebookClient 封装 Notebook CRD 操作
type NotebookClient struct {
	client ctrlclient.Client
}

func NewNotebookClient(client ctrlclient.Client) *NotebookClient {
	return &NotebookClient{client: client}
}

func (c *NotebookClient) Create(ctx context.Context, notebook *nbv1.Notebook) error {
	return c.client.Create(ctx, notebook)
}

func (c *NotebookClient) Get(ctx context.Context, namespace, name string) (*nbv1.Notebook, error) {
	notebook := &nbv1.Notebook{}
	err := c.client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, notebook)
	return notebook, err
}

func (c *NotebookClient) Update(ctx context.Context, notebook *nbv1.Notebook) error {
	return c.client.Update(ctx, notebook)
}

func (c *NotebookClient) Delete(ctx context.Context, notebook *nbv1.Notebook, opts ...ctrlclient.DeleteOption) error {
	return c.client.Delete(ctx, notebook, opts...)
}

func (c *NotebookClient) List(ctx context.Context, namespace string, opts ...ctrlclient.ListOption) (*nbv1.NotebookList, error) {
	list := &nbv1.NotebookList{}
	opts = append(opts, ctrlclient.InNamespace(namespace))
	err := c.client.List(ctx, list, opts...)
	return list, err
}

// TensorboardClient 封装 Tensorboard CRD 操作
type TensorboardClient struct {
	client ctrlclient.Client
}

func NewTensorboardClient(client ctrlclient.Client) *TensorboardClient {
	return &TensorboardClient{client: client}
}

func (c *TensorboardClient) Create(ctx context.Context, tb *tbv1alpha1.Tensorboard) error {
	return c.client.Create(ctx, tb)
}

func (c *TensorboardClient) Get(ctx context.Context, namespace, name string) (*tbv1alpha1.Tensorboard, error) {
	tb := &tbv1alpha1.Tensorboard{}
	err := c.client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, tb)
	return tb, err
}

func (c *TensorboardClient) Update(ctx context.Context, tb *tbv1alpha1.Tensorboard) error {
	return c.client.Update(ctx, tb)
}

func (c *TensorboardClient) Delete(ctx context.Context, tb *tbv1alpha1.Tensorboard, opts ...ctrlclient.DeleteOption) error {
	return c.client.Delete(ctx, tb, opts...)
}

func (c *TensorboardClient) List(ctx context.Context, namespace string, opts ...ctrlclient.ListOption) (*tbv1alpha1.TensorboardList, error) {
	list := &tbv1alpha1.TensorboardList{}
	opts = append(opts, ctrlclient.InNamespace(namespace))
	err := c.client.List(ctx, list, opts...)
	return list, err
}
