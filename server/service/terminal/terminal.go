package terminal

import (
	"context"
	"gin-vue-admin/global"
	"io"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// TerminalService 提供通用的 Web Terminal 服务
// 可被 Training、Notebook 等多种资源使用
type TerminalService struct{}

var TerminalServiceApp = new(TerminalService)

// TerminalRequest 终端连接请求
type TerminalRequest struct {
	Namespace     string   // Pod 所在命名空间
	PodName       string   // Pod 名称（直接指定）
	LabelSelector string   // 或通过标签选择器查找 Pod
	Container     string   // 容器名称（可选）
	Command       []string // 执行的命令（默认 /bin/sh）
	ClusterID     uint     // 集群 ID
}

// HandleWebSocket 处理 WebSocket 连接并执行 kubectl exec
func (s *TerminalService) HandleWebSocket(ctx context.Context, req *TerminalRequest, conn *websocket.Conn) error {
	// 1. 获取集群客户端
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(req.ClusterID)
	if cluster == nil {
		return errors.New("集群不存在")
	}

	// 2. 确定 Pod 名称
	podName := req.PodName
	if podName == "" && req.LabelSelector != "" {
		var err error
		podName, err = s.findRunningPod(ctx, cluster.ClientSet, req.Namespace, req.LabelSelector)
		if err != nil {
			return err
		}
	}
	if podName == "" {
		return errors.New("未指定 Pod")
	}

	// 3. 确定命令
	command := req.Command
	if len(command) == 0 {
		command = []string{"/bin/sh"}
	}

	// 4. 执行 exec
	return s.executeExec(ctx, cluster.RestConfig, cluster.ClientSet, req.Namespace, podName, req.Container, command, conn)
}

// findRunningPod 根据标签选择器查找运行中的 Pod
func (s *TerminalService) findRunningPod(ctx context.Context, clientSet *kubernetes.Clientset, namespace, labelSelector string) (string, error) {
	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return "", errors.Wrap(err, "查找 Pod 失败")
	}

	if len(pods.Items) == 0 {
		return "", errors.New("未找到 Pod")
	}

	// 优先返回运行中的 Pod
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			return pod.Name, nil
		}
	}

	return pods.Items[0].Name, nil
}

// executeExec 执行 kubectl exec
func (s *TerminalService) executeExec(
	ctx context.Context,
	restConfig *rest.Config,
	clientSet *kubernetes.Clientset,
	namespace, podName, containerName string,
	command []string,
	conn *websocket.Conn,
) error {
	// 构建 exec 请求
	req := clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	execOptions := &corev1.PodExecOptions{
		Command: command,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	if containerName != "" {
		execOptions.Container = containerName
	}

	req.VersionedParams(execOptions, scheme.ParameterCodec)

	// 创建 SPDY 执行器
	exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
	if err != nil {
		return errors.Wrap(err, "创建SPDY执行器失败")
	}

	// 创建 WebSocket 流适配器
	streamHandler := &wsStreamHandler{conn: conn}

	logx.Info("开始 Terminal 会话",
		logx.Field("namespace", namespace),
		logx.Field("pod", podName),
		logx.Field("container", containerName),
	)

	// 执行
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  streamHandler,
		Stdout: streamHandler,
		Stderr: streamHandler,
		Tty:    true,
	})

	if err != nil {
		logx.Error("Terminal 会话异常结束", logx.Field("err", err))
	} else {
		logx.Info("Terminal 会话正常结束")
	}

	return err
}

// GetPodList 获取可连接终端的 Pod 列表
func (s *TerminalService) GetPodList(ctx context.Context, clusterID uint, namespace, labelSelector string) ([]PodInfo, error) {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterID)
	if cluster == nil {
		return nil, errors.New("集群不存在")
	}

	pods, err := cluster.ClientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, errors.Wrap(err, "获取 Pod 列表失败")
	}

	result := make([]PodInfo, 0, len(pods.Items))
	for _, pod := range pods.Items {
		info := PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			HostIP:    pod.Status.HostIP,
			PodIP:     pod.Status.PodIP,
		}

		// 获取容器列表
		for _, container := range pod.Spec.Containers {
			info.Containers = append(info.Containers, container.Name)
		}

		result = append(result, info)
	}

	return result, nil
}

// PodInfo Pod 信息
type PodInfo struct {
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	Status     string   `json:"status"`
	HostIP     string   `json:"hostIP"`
	PodIP      string   `json:"podIP"`
	Containers []string `json:"containers"`
}

// wsStreamHandler 实现 io.Reader, io.Writer 适配 WebSocket
type wsStreamHandler struct {
	conn   *websocket.Conn
	buffer []byte
}

func (h *wsStreamHandler) Read(p []byte) (n int, err error) {
	if len(h.buffer) > 0 {
		n = copy(p, h.buffer)
		h.buffer = h.buffer[n:]
		return n, nil
	}

	_, message, err := h.conn.ReadMessage()
	if err != nil {
		// 如果前端触发了 ws.close() 或是网络正常断开，把它当作为干净的 IO EOF 流结束
		// 以免 K8s 的 SPDY Executor 惊慌失措并报 Unhandled Error，防止遗留僵尸 exec 进程
		if websocket.IsUnexpectedCloseError(err) || websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
			return 0, io.EOF
		}
		// 其它常规断开错误也统一以 EOF 退出
		return 0, io.EOF
	}

	n = copy(p, message)
	if n < len(message) {
		h.buffer = message[n:]
	}
	return n, nil
}

func (h *wsStreamHandler) Write(p []byte) (n int, err error) {
	err = h.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
