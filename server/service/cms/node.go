package cms

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	clusterModel "gin-vue-admin/model/cluster"
	cmsRes "gin-vue-admin/model/cms/response"
	"gin-vue-admin/model/consts"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeService struct{}

// GetClusterNodes 获取集群下的所有Node节点信息（从内存集群管理器获取集群元数据）
func (s *NodeService) GetClusterNodes(clusterId uint, keyword string) ([]cmsRes.NodeInfoResponse, error) {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterId)
	if cluster == nil {
		return nil, errors.Errorf("集群不存在或未连接")
	}
	if cluster.ClientSet == nil {
		return nil, errors.Errorf("集群K8s客户端未初始化")
	}

	nodes, err := s.ListNodeResources(context.Background(), cluster.ClientSet, keyword)
	if err != nil {
		return nil, err
	}

	for i := range nodes {
		nodes[i].Area = cluster.Area
		nodes[i].ClusterName = cluster.ClusterName
	}
	return nodes, nil
}

// GetClusterNodesWithResources 获取集群节点及详细资源信息（从DB获取集群元数据）
func (s *NodeService) GetClusterNodesWithResources(ctx context.Context, clusterId uint, keyword string) ([]cmsRes.NodeInfoResponse, error) {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterId)
	if cluster == nil {
		return nil, errors.Errorf("集群不存在或未连接")
	}
	if cluster.ClientSet == nil {
		return nil, errors.Errorf("集群K8s客户端未初始化")
	}

	nodes, err := s.ListNodeResources(ctx, cluster.ClientSet, keyword)
	if err != nil {
		return nil, err
	}

	var clusterInfo clusterModel.K8sCluster
	if err := global.GVA_DB.Model(&clusterModel.K8sCluster{}).Where("id = ?", clusterId).First(&clusterInfo).Error; err != nil {
		return nil, errors.Errorf("查询集群信息失败: %v", err)
	}

	for i := range nodes {
		nodes[i].Area = clusterInfo.Area
		nodes[i].ClusterName = clusterInfo.Name
	}
	return nodes, nil
}

// ---------------------- 节点操作 ----------------------

// UncordonNode 恢复节点调度
func (s *NodeService) UncordonNode(ctx context.Context, clusterId uint, nodeName string) error {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterId)
	if cluster == nil {
		return errors.Errorf("集群不存在或未连接")
	}

	node, err := cluster.ClientSet.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return errors.Errorf("查询节点失败: %v", err)
	}

	if !node.Spec.Unschedulable {
		return nil
	}

	node.Spec.Unschedulable = false
	_, err = cluster.ClientSet.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
	if err != nil {
		return errors.Errorf("更新节点调度状态失败: %v", err)
	}
	return nil
}

// DrainNode 驱逐节点 (Cordon + Evict)
func (s *NodeService) DrainNode(ctx context.Context, clusterId uint, nodeName string) error {
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(clusterId)
	if cluster == nil {
		return errors.Errorf("集群不存在或未连接")
	}

	node, err := cluster.ClientSet.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return errors.Errorf("查询节点失败: %v", err)
	}

	if !node.Spec.Unschedulable {
		node.Spec.Unschedulable = true
		_, err = cluster.ClientSet.CoreV1().Nodes().Update(ctx, node, metav1.UpdateOptions{})
		if err != nil {
			return errors.Errorf("停止节点调度失败: %v", err)
		}
	}

	return evictPodsOnNode(ctx, cluster.ClientSet, nodeName)
}

// evictPodsOnNode 驱逐指定节点上的所有非 DaemonSet Pod
func evictPodsOnNode(ctx context.Context, clientSet *kubernetes.Clientset, nodeName string) error {
	pods, err := clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		return errors.Errorf("获取节点上 Pod 列表失败: %v", err)
	}

	for _, pod := range pods.Items {
		isDaemonSet := false
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == "DaemonSet" {
				isDaemonSet = true
				break
			}
		}
		if isDaemonSet {
			continue
		}

		gracePeriod := int64(30)
		err = clientSet.CoreV1().Pods(pod.Namespace).Delete(ctx, pod.Name, metav1.DeleteOptions{
			GracePeriodSeconds: &gracePeriod,
		})
		if err != nil {
			logx.Errorf("驱逐 Pod %s/%s 失败: %v", pod.Namespace, pod.Name, err)
			continue
		}
	}
	return nil
}

// CalculateNodeResources 获取单个节点的资源信息（Allocatable、GPU、vGPU、CPU型号）
func (s *NodeService) CalculateNodeResources(ctx context.Context, clientSet *kubernetes.Clientset, nodeName string) (*cmsRes.NodeInfoResponse, error) {
	node, err := clientSet.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Errorf("获取节点信息失败: %v", err)
	}

	info := &cmsRes.NodeInfoResponse{
		NodeName:    node.Name,
		Schedulable: !node.Spec.Unschedulable,
		NodeRole:    parseNodeRole(node),
	}

	// 内网 IP
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			info.InternalIP = addr.Address
			break
		}
	}

	// CPU
	if cpuAlloc, ok := node.Status.Allocatable[corev1.ResourceCPU]; ok {
		info.CPUAllocatable = cpuAlloc.Value()
		info.CPUAvailable = cpuAlloc.Value()
	}
	if cpuModel, ok := node.Labels[consts.LabelCPUModel]; ok {
		info.CPUModel = cpuModel
	}

	// 内存
	if memAlloc, ok := node.Status.Allocatable[corev1.ResourceMemory]; ok {
		memGB := memAlloc.Value() / (1024 * 1024 * 1024)
		info.MemoryAllocatable = memGB
		info.MemoryAvailable = memGB
	}

	// GPU / vGPU
	parseGPUInfo(node, info)
	info.GPUAvailable = info.GPUCount

	return info, nil
}

// CalculateMaxInstances 根据节点资源和产品规格计算最大实例数（GPU / vGPU / CPU-only 互斥）
func (s *NodeService) CalculateMaxInstances(nodeInfo *cmsRes.NodeInfoResponse, gpuCount, vgpuNumber, vgpuMemory, vgpuCores, cpu, memory int64, isGPU, isVGPU bool) int64 {
	switch {
	case isGPU:
		if gpuCount > 0 {
			return safeDiv(nodeInfo.GPUCount, gpuCount)
		}
		return 0
	case isVGPU:
		return minByRequested(
			resourceQuota{total: nodeInfo.VGPUNumber, requested: vgpuNumber},
			resourceQuota{total: nodeInfo.VGPUMemory, requested: vgpuMemory},
			resourceQuota{total: nodeInfo.VGPUCores, requested: vgpuCores},
		)
	default:
		return minByRequested(
			resourceQuota{total: nodeInfo.CPUAllocatable, requested: cpu},
			resourceQuota{total: nodeInfo.MemoryAllocatable, requested: memory},
		)
	}
}

// ListNodeResources 获取集群所有节点的资源信息（支持关键字过滤）
func (s *NodeService) ListNodeResources(ctx context.Context, clientSet *kubernetes.Clientset, keyword string) ([]cmsRes.NodeInfoResponse, error) {
	nodeList, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Errorf("获取节点列表失败: %v", err)
	}

	calc := &NodeService{}
	result := make([]cmsRes.NodeInfoResponse, 0, len(nodeList.Items))
	for _, node := range nodeList.Items {
		if keyword != "" && !nodeMatchesKeyword(&node, keyword) {
			continue
		}
		info, err := calc.CalculateNodeResources(ctx, clientSet, node.Name)
		if err != nil {
			logx.Errorf("获取节点 %s 资源失败: %v", node.Name, err)
			continue
		}
		result = append(result, *info)
	}
	return result, nil
}

// nodeMatchesKeyword 判断节点名称或IP是否匹配关键字
func nodeMatchesKeyword(node *corev1.Node, keyword string) bool {
	if strings.Contains(node.Name, keyword) {
		return true
	}
	for _, addr := range node.Status.Addresses {
		if strings.Contains(addr.Address, keyword) {
			return true
		}
	}
	return false
}

// parseGPUInfo 解析 GPU 和 vGPU 信息
func parseGPUInfo(node *corev1.Node, info *cmsRes.NodeInfoResponse) {
	labels := node.Labels

	// GPU 型号
	if v, ok := labels[consts.LabelGpuProduct]; ok {
		info.GPUModel = v
	} else if v, ok := labels[consts.LabelGPUModel]; ok {
		info.GPUModel = v
	} else if v, ok := labels[consts.LabelAccelerator]; ok {
		info.GPUModel = v
	}

	// GPU 数量
	if gpu, ok := node.Status.Capacity[corev1.ResourceName(consts.NvidiaGPUType)]; ok {
		info.GPUCount = gpu.Value()
	} else if gpu, ok := node.Status.Capacity[corev1.ResourceName(consts.AmdGPUType)]; ok {
		info.GPUCount = gpu.Value()
	}

	// GPU 显存（标签值为 MB，转 GB）
	if memStr, ok := labels[consts.LabelGPUMem]; ok {
		if mem, err := parseGPUMemoryMB(memStr); err == nil {
			info.GPUMemory = mem / 1024
		}
	}

	// vGPU 资源（Volcano vGPU 调度器）
	if v, ok := node.Status.Capacity[corev1.ResourceName(consts.VolcanoVGPUNumber)]; ok {
		info.VGPUNumber = v.Value()
	}
	if v, ok := node.Status.Capacity[corev1.ResourceName(consts.VolcanoVGPUMemory)]; ok {
		info.VGPUMemory = v.Value()
	}
	if v, ok := node.Status.Capacity[corev1.ResourceName(consts.VolcanoVGPUCores)]; ok {
		info.VGPUCores = v.Value()
	}
}

// parseNodeRole 根据节点标签判断角色
func parseNodeRole(node *corev1.Node) string {
	if _, ok := node.Labels[consts.LabelNodeRoleMaster]; ok {
		return "master"
	}
	if _, ok := node.Labels[consts.LabelNodeRoleControlPlane]; ok {
		return "master"
	}
	return "worker"
}

// parseGPUMemoryMB 解析 GPU 显存字符串（MB）
func parseGPUMemoryMB(memStr string) (int64, error) {
	memStr = strings.TrimSpace(memStr)
	memStr = strings.TrimSuffix(memStr, "MiB")
	memStr = strings.TrimSuffix(memStr, "MB")
	memStr = strings.TrimSuffix(memStr, "Mi")
	memStr = strings.TrimSpace(memStr)

	var mem int64
	if _, err := fmt.Sscanf(memStr, "%d", &mem); err != nil {
		return 0, err
	}
	return mem, nil
}

// safeDiv 安全除法，除数为0时返回0
func safeDiv(total, per int64) int64 {
	if per <= 0 || total <= 0 {
		return 0
	}
	return total / per
}

// resourceQuota 表示一种资源的总量和请求量
type resourceQuota struct {
	total     int64
	requested int64
}

// minByRequested 计算多种资源约束下的最大实例数
// 规则：requested > 0 表示该资源被请求，必须参与计算（即使结果为0也是硬约束）
//
//	requested <= 0 表示未请求该资源，跳过
func minByRequested(quotas ...resourceQuota) int64 {
	result := int64(-1) // -1 表示尚未计算
	for _, q := range quotas {
		if q.requested <= 0 {
			continue // 未请求该资源，跳过
		}
		n := safeDiv(q.total, q.requested)
		if result < 0 || n < result {
			result = n
		}
	}
	if result < 0 {
		return 0
	}
	return result
}
