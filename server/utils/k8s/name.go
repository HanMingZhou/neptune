package helper

import (
	"context"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RFC 1123 子域名正则（K8s 资源名称规范）
var rfc1123Regex = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)

// ValidateNotebookName 校验 Notebook 名称是否符合 K8s 命名规范
func ValidateNotebookName(name string) error {
	if name == "" {
		return errors.New("Notebook名称不能为空")
	}
	if len(name) > 63 {
		return errors.New("Notebook名称不能超过63个字符")
	}
	if !rfc1123Regex.MatchString(name) {
		return errors.New("Notebook名称必须由小写字母、数字、'-'或'.'组成，且以字母或数字开头和结尾（例如：my-notebook）")
	}
	return nil
}

// 根据实例类型生成对应的实例Name: instanceType+"-"+6位uuid
func GenerateInstanceName(instanceType string) string {
	// 生成 UUID 并取前 6 位（去掉 '-' 后取前 6 个字符更短、更美观）
	u := strings.ReplaceAll(uuid.New().String(), "-", "")
	shortUUID := u[:6]
	return string(instanceType) + "-" + shortUUID
}

// 生成namespace name
func GenerateNamespaceName(uuid string) string {
	// 生成 UUID 并取前 6 位（去掉 '-' 后取前 6 个字符更短、更美观）
	u := strings.ReplaceAll(uuid, "-", "")
	shortUUID := u[:6]
	return shortUUID
}

// EnsureNamespace 确保 namespace 存在
func EnsureNamespace(ctx context.Context, clientSet kubernetes.Interface, namespace string) error {
	_, err := clientSet.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err == nil {
		return nil
	}

	// 创建 namespace
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	if _, err := clientSet.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{}); err != nil {
		logx.Error("创建namespace失败", logx.Field("err", err))
		return errors.Wrap(err, "创建命名空间失败")
	}
	return nil
}
