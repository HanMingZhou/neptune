package secret

import (
	"context"
	"fmt"
	secretModel "gin-vue-admin/model/secret/request"

	"github.com/zeromicro/go-zero/core/logx"
	corev1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SecretManager interface {
	CreateSSHSecret(ctx context.Context, req *secretModel.AddSecretReq) error
	DeleteSSHSecret(ctx context.Context, req *secretModel.DeleteSecretReq) error
	CreateSSHPrivateKeySecret(ctx context.Context, req *secretModel.AddSecretReq) (string, error)
	DeleteSSHPrivateKeySecret(ctx context.Context, req *secretModel.DeleteSecretReq) error
	CreateSSHPasswordSecret(ctx context.Context, req *secretModel.AddSecretReq) error
	DeleteSSHPasswordSecret(ctx context.Context, req *secretModel.DeleteSecretReq) error
}

var _ SecretManager = (*K8sSecretService)(nil)

type K8sSecretService struct {
	client kubernetes.Interface
}

func NewK8sSecretManager(client kubernetes.Interface) *K8sSecretService {
	return &K8sSecretService{client: client}
}

func (m *K8sSecretService) CreateSSHSecret(ctx context.Context, req *secretModel.AddSecretReq) error {
	if req.Content == nil {
		return nil
	}

	secretName := fmt.Sprintf("%s-ssh-key", req.InstanceName)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"app":            req.InstanceName,
				req.InstanceType: req.InstanceName,
			},
		},
		StringData: req.Content,
		Type:       corev1.SecretTypeOpaque, // 或者 kubernetes.io/ssh-auth
	}

	if _, err := m.client.CoreV1().
		Secrets(req.Namespace).
		Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		logx.Error("创建SSH Secret失败", err)
		return err
	}
	return nil
}

func (m *K8sSecretService) DeleteSSHSecret(ctx context.Context, req *secretModel.DeleteSecretReq) error {
	secretName := fmt.Sprintf("%s-ssh-key", req.InstanceName)
	if err := m.client.CoreV1().
		Secrets(req.Namespace).
		Delete(ctx, secretName, metav1.DeleteOptions{}); err != nil {
		if k8serr.IsNotFound(err) {
			return nil
		}
		logx.Error("删除SSH Secret失败", err)
		return err
	}
	return nil
}

func (m *K8sSecretService) CreateSSHPrivateKeySecret(ctx context.Context, req *secretModel.AddSecretReq) (string, error) {
	secretName := fmt.Sprintf("%s-ssh-private-key", req.InstanceName)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"app":            req.InstanceName,
				req.InstanceType: req.InstanceName,
			},
		},
		StringData: req.Content,
		Type:       corev1.SecretTypeOpaque,
	}

	if _, err := m.client.CoreV1().
		Secrets(req.Namespace).
		Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		logx.Error("创建SSH私钥Secret失败", err)
		return "", err
	}
	return secretName, nil
}

func (m *K8sSecretService) DeleteSSHPrivateKeySecret(ctx context.Context, req *secretModel.DeleteSecretReq) error {
	secretName := fmt.Sprintf("%s-ssh-private-key", req.InstanceName)
	if err := m.client.CoreV1().
		Secrets(req.Namespace).
		Delete(ctx, secretName, metav1.DeleteOptions{}); err != nil {
		if k8serr.IsNotFound(err) {
			return nil
		}
		logx.Error("删除SSH私钥Secret失败", err)
		return err
	}
	return nil
}

// CreateSSHPasswordSecret 创建SSH密码Secret
func (m *K8sSecretService) CreateSSHPasswordSecret(ctx context.Context, req *secretModel.AddSecretReq) error {
	if req.Content == nil {
		return nil
	}

	secretName := fmt.Sprintf("%s-ssh-password", req.InstanceName)

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"app":            req.InstanceName,
				req.InstanceType: req.InstanceName,
			},
		},
		StringData: req.Content,
		Type:       corev1.SecretTypeOpaque,
	}

	if _, err := m.client.CoreV1().
		Secrets(req.Namespace).
		Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		logx.Error("创建SSH密码Secret失败", err)
		return err
	}
	return nil
}

// DeleteSSHPasswordSecret 删除SSH密码Secret
func (m *K8sSecretService) DeleteSSHPasswordSecret(ctx context.Context, req *secretModel.DeleteSecretReq) error {
	secretName := fmt.Sprintf("%s-ssh-password", req.InstanceName)
	if err := m.client.CoreV1().
		Secrets(req.Namespace).
		Delete(ctx, secretName, metav1.DeleteOptions{}); err != nil {
		if k8serr.IsNotFound(err) {
			return nil
		}
		logx.Error("删除SSH密码Secret失败", err)
		return err
	}
	return nil
}
