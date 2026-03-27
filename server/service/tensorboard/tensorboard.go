package tensorboard

import (
	"context"
	"gin-vue-admin/global"
	tensorboardModel "gin-vue-admin/model/tensorboard/request"

	tbv1alpha1 "github.com/kubeflow/kubeflow/components/tensorboard-controller/api/v1alpha1"
	"github.com/zeromicro/go-zero/core/logx"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TensorboardManager interface {
	CreateTensorboard(ctx context.Context, req *tensorboardModel.AddTensorBoardReq) error
	DeleteTensorboard(ctx context.Context, req *tensorboardModel.DeleteTensorBoardReq) error
}

var _ TensorboardManager = &K8sTensorboardService{}

type K8sTensorboardService struct {
	client *global.TensorboardClient
}

func NewTensorboardManager(client *global.TensorboardClient) *K8sTensorboardService {
	return &K8sTensorboardService{client: client}
}

func (m *K8sTensorboardService) CreateTensorboard(ctx context.Context, req *tensorboardModel.AddTensorBoardReq) error {
	if !req.EnableTensorboard {
		return nil
	}

	// 确定 LogDir
	// 注意：Tensorboard 需要访问持久化存储
	// 这里假设 workspace 卷是一个 PVC，且日志存储在 logs 目录下
	// 格式: pvc://<pvc-name>/<sub-path>

	tb := &tbv1alpha1.Tensorboard{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.InstanceName,
			Namespace: req.Namespace,
			Labels:    req.Labels,
		},
		Spec: tbv1alpha1.TensorboardSpec{
			LogsPath: req.LogsPath,
		},
	}

	if err := m.client.Create(ctx, tb); err != nil {
		logx.Error("创建Tensorboard失败", err)
		return err
	}
	return nil
}

func (m *K8sTensorboardService) DeleteTensorboard(ctx context.Context, req *tensorboardModel.DeleteTensorBoardReq) error {
	tb := &tbv1alpha1.Tensorboard{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.InstanceName,
			Namespace: req.Namespace,
		},
	}
	if err := m.client.Delete(ctx, tb); err != nil {
		logx.Error("删除Tensorboard失败", err)
		return err
	}
	return nil
}
