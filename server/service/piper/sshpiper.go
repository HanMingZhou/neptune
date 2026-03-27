package piper

import (
	"context"
	"fmt"
	pipeModel "gin-vue-admin/model/pipe/request"
	"strings"

	piperv1beta1 "github.com/tg123/sshpiper/plugin/kubernetes/apis/sshpiper/v1beta1"
	"github.com/zeromicro/go-zero/core/logx"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type SSHPiperManager interface {
	CreatePipe(ctx context.Context, req *pipeModel.AddPipeReq) error
	DeletePipe(ctx context.Context, req *pipeModel.DeletePipeReq) error
}

var _ SSHPiperManager = &K8sSSHPiperService{}

type K8sSSHPiperService struct {
	client ctrlclient.Client
}

func NewK8sSSHPiperManager(client ctrlclient.Client) *K8sSSHPiperService {
	return &K8sSSHPiperService{client: client}
}

func (m *K8sSSHPiperService) CreatePipe(ctx context.Context, req *pipeModel.AddPipeReq) error {
	// 如果既没有公钥也没有启用密码登录，则不创建 Pipe
	if req.UserSSHKey == "" && !req.EnablePasswordAuth {
		return nil
	}

	// 构造 Pipe 名称: namespace-name
	// 确定目标命名空间
	targetNamespace := req.TargetNamespace
	if targetNamespace == "" {
		targetNamespace = req.Namespace
	}

	// 构造 Pipe 名称: pipe-namespace-instanceName
	// 例如: pipe-hmz-notebook-bc276e
	pipeName := fmt.Sprintf("pipe-%s-%s", targetNamespace, req.InstanceName)

	// 构造用户名 (用于 SSH 登录): targetNamespace-instanceName
	username := fmt.Sprintf("%s-%s", targetNamespace, req.InstanceName)

	// 构造后端地址
	// Notebook Service 地址: <notebook-name>.<target-namespace>.svc.cluster.local:22
	targetHost := req.TargetHost
	if targetHost == "" {
		targetHost = fmt.Sprintf("%s.%s.svc.cluster.local:22", req.InstanceName, targetNamespace)
	}

	// 构造 FromSpec
	fromSpec := piperv1beta1.FromSpec{
		Username: username,
	}

	// 公钥认证：设置 AuthorizedKeysData
	// SSHPiper 会验证用户的公钥
	if req.UserSSHKey != "" {
		fromSpec.AuthorizedKeysData = strings.TrimSpace(req.UserSSHKey)
	}

	// 密码认证：不设置 HtpasswdData
	// SSHPiper 会将密码透传给后端 Pod，由 Pod 的 sshd 验证
	// 这样用户在容器内执行 passwd 后，新密码立即生效

	// 确定后端用户名（默认 root）
	targetUsername := req.TargetUsername
	if targetUsername == "" {
		targetUsername = "root" // 默认使用 root 用户
	}

	// 构造 ToSpec
	toSpec := piperv1beta1.ToSpec{
		Host:          targetHost,
		Username:      targetUsername, // 后端容器内的用户名（可配置）
		IgnoreHostkey: true,           // 忽略后端 HostKey 检查
	}

	// 设置后端认证方式
	// 公钥登录：SSHPiper 用私钥连接后端
	// 密码登录：不设置凭证，SSHPiper 透传密码给后端
	if req.PrivateKeySecretName != "" {
		toSpec.PrivateKeySecret = corev1.LocalObjectReference{
			Name: req.PrivateKeySecretName,
		}
	}
	// 密码登录时不设置 PrivateKeySecret 或 PasswordSecret
	// SSHPiper 会将用户输入的密码透传给后端 Pod

	pipe := &piperv1beta1.Pipe{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pipeName,
			Namespace: req.Namespace,
			Labels:    req.Labels,
		},
		Spec: piperv1beta1.PipeSpec{
			From: []piperv1beta1.FromSpec{fromSpec},
			To:   toSpec,
		},
	}

	if err := m.client.Create(ctx, pipe); err != nil {
		logx.Error("创建SSH Pipe失败", err)
		return err
	}
	return nil
}

func (m *K8sSSHPiperService) DeletePipe(ctx context.Context, req *pipeModel.DeletePipeReq) error {
	// Pipe 名称格式: <target-namespace>-<instance-name>
	// req.Namespace: Pipe 所在的 namespace (通常是 kubeflow)
	// req.InstanceName: Pipe 名称（已经是完整名称，例如 hmz-notebook-879264）
	pipe := &piperv1beta1.Pipe{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.InstanceName, // 直接使用完整名称 (调用方应已拼接好 pipe- 前缀)
			Namespace: req.Namespace,    // Pipe 所在 namespace
		},
	}
	if err := m.client.Delete(ctx, pipe); err != nil {
		if ctrlclient.IgnoreNotFound(err) != nil {
			logx.Error("删除SSH Pipe失败", err)
			return err
		}
	}
	return nil
}
