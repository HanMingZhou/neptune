package global

import (
	"context"
	"fmt"
	clusterModel "gin-vue-admin/model/cluster"
	"gin-vue-admin/service/podgroup"
	"sync"
	"time"

	apisixclient "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned"
	nbv1 "github.com/kubeflow/kubeflow/components/notebook-controller/api/v1"
	tbv1alpha1 "github.com/kubeflow/kubeflow/components/tensorboard-controller/api/v1alpha1"
	piperv1beta1 "github.com/tg123/sshpiper/plugin/kubernetes/apis/sshpiper/v1beta1"
	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1lister "k8s.io/client-go/listers/apps/v1"
	v1podlister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	vcinformers "volcano.sh/apis/pkg/client/informers/externalversions"
)

type ClusterManager interface {
	AddCluster(clusterId uint, kubeConfig []byte, clusterName, area string) error
	GetCluster(clusterId uint) *ClusterClientInfo
	RemoveCluster(clusterId uint) error
	ReloadCluster(clusterId uint, kubeConfig []byte, clusterName, area string) error
	ListClusters() []*ClusterClientInfo
	GetClusterCount() int
	Shutdown()
}

// ClusterClientInfo 单个集群的客户端信息
type ClusterClientInfo struct {
	ClusterId           uint                              // 集群ID
	ClusterName         string                            // 集群名称
	Area                string                            // 地域
	DefaultStorageClass string                            // 默认 StorageClass
	ClientSet           *kubernetes.Clientset             // 标准客户端
	DynamicClient       dynamic.Interface                 // 动态客户端
	RestConfig          *rest.Config                      // REST配置
	RuntimeClient       ctrlclient.Client                 // controller-runtime客户端
	K8sInformer         informers.SharedInformerFactory   // Informer工厂
	PodLister           v1podlister.PodLister             // Pod Lister
	StsLister           v1lister.StatefulSetLister        // StatefulSet Lister
	DeployLister        v1lister.DeploymentLister         // Deployment Lister
	NotebookClient      *NotebookClient                   // Notebook客户端
	TensorboardClient   *TensorboardClient                // Tensorboard客户端
	ApisixClient        apisixclient.Interface            // Apisix客户端
	VolcanoClient       *VolcanoClient                    // Volcano客户端
	PodGroupInformer    *podgroup.PodGroupInformerFactory // PodGroup Informer
	StopCh              chan struct{}                     // 停止通道
}

// K8sClusterManager 实现ClusterManager接口
var _ ClusterManager = (*K8sClusterManager)(nil)

// K8sClusterManager K8s集群管理器
type K8sClusterManager struct {
	mu           sync.RWMutex
	clusters     map[uint]*ClusterClientInfo // key: clusterID
	orderManager podgroup.OrderManager
}

// 注意：GVA_K8S_CLUSTER_MANAGER 全局变量在 global.go 中声明

// NewK8sClusterManager 创建集群管理器
func NewK8sClusterManager() *K8sClusterManager {
	return &K8sClusterManager{
		clusters: make(map[uint]*ClusterClientInfo),
	}
}

// SetOrderManager 设置计费管理器，并同步更新已创建的集群
func (m *K8sClusterManager) SetOrderManager(bm podgroup.OrderManager) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.orderManager = bm
	for _, client := range m.clusters {
		if client.PodGroupInformer != nil {
			client.PodGroupInformer.SetOrderManager(bm)
		}
	}
}

// SetProductManager 设置产品管理器，并同步更新已创建的集群
func (m *K8sClusterManager) SetProductManager(pm podgroup.ProductManager) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, client := range m.clusters {
		if client.PodGroupInformer != nil {
			client.PodGroupInformer.SetProductManager(pm)
		}
	}
}

// AddCluster 添加集群到管理器
func (m *K8sClusterManager) AddCluster(clusterId uint, kubeConfig []byte, clusterName, area string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在
	if _, exists := m.clusters[clusterId]; exists {
		return fmt.Errorf("集群 %d 已存在", clusterId)
	}

	// 从kubeconfig创建client
	clientInfo, err := createClusterClient(clusterId, kubeConfig, clusterName, area, m.orderManager)
	if err != nil {
		return fmt.Errorf("创建集群 %d 客户端失败: %w", clusterId, err)
	}

	m.clusters[clusterId] = clientInfo

	logx.Info("添加集群成功", clusterId, clusterName)
	return nil
}

// GetCluster 根据clusterID获取集群客户端
func (m *K8sClusterManager) GetCluster(clusterId uint) *ClusterClientInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clusters[clusterId]
	if !exists {
		m.mu.RUnlock() // 必须先释放读锁，因为 AddCluster 需要获取写锁

		// 如果不存在，则查询数据库配置后创建一个加入到管理器中
		var clusterConfig clusterModel.K8sCluster
		GVA_DB.Where("id = ?", clusterId).First(&clusterConfig)
		if clusterConfig.ID == 0 {
			return nil
		}
		m.AddCluster(clusterId, []byte(clusterConfig.KubeConfig), clusterConfig.Name, clusterConfig.Area)

		m.mu.RLock() // 重新获取读锁，以便读取刚才添加的 client，并满足 defer 的解锁要求
		client = m.clusters[clusterId]
	}
	return client
}

// RemoveCluster 移除集群
func (m *K8sClusterManager) RemoveCluster(clusterId uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.clusters[clusterId]
	if !exists {
		return fmt.Errorf("集群 %d 不存在", clusterId)
	}

	// 停止informer
	close(client.StopCh)
	if client.PodGroupInformer != nil {
		client.PodGroupInformer.Stop()
	}
	delete(m.clusters, clusterId)
	return nil
}

// ReloadCluster 重载集群：销毁旧连接，用新配置重建客户端
func (m *K8sClusterManager) ReloadCluster(clusterId uint, kubeConfig []byte, clusterName, area string) error {
	m.mu.Lock()

	// 先清理旧连接（如果存在）
	if old, exists := m.clusters[clusterId]; exists {
		close(old.StopCh)
		if old.PodGroupInformer != nil {
			old.PodGroupInformer.Stop()
		}
		delete(m.clusters, clusterId)
	}

	m.mu.Unlock()

	// 重新创建并添加
	if err := m.AddCluster(clusterId, kubeConfig, clusterName, area); err != nil {
		return fmt.Errorf("重载集群 %d 失败: %w", clusterId, err)
	}

	logx.Info("重载集群成功", logx.Field("clusterId", clusterId), logx.Field("clusterName", clusterName))
	return nil
}

// ListClusters 列出所有集群
func (m *K8sClusterManager) ListClusters() []*ClusterClientInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clusters := make([]*ClusterClientInfo, 0, len(m.clusters))
	for _, client := range m.clusters {
		clusters = append(clusters, client)
	}
	return clusters
}

// GetClusterCount 获取集群数量
func (m *K8sClusterManager) GetClusterCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clusters)
}

// GetClusterByArea 根据区域获取集群（返回第一个匹配的）
func (m *K8sClusterManager) GetClusterByArea(area string) *ClusterClientInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, client := range m.clusters {
		if client.Area == area {
			return client
		}
	}

	// 如果内存中没有，尝试从数据库加载
	var clusterConfig clusterModel.K8sCluster
	GVA_DB.Where("area = ?", area).First(&clusterConfig)
	if clusterConfig.ID == 0 {
		return nil
	}

	// 加载到管理器
	_ = m.AddCluster(clusterConfig.ID, []byte(clusterConfig.KubeConfig), clusterConfig.Name, clusterConfig.Area)
	return m.clusters[clusterConfig.ID]
}

// GetAllAreas 获取所有可用区域
func (m *K8sClusterManager) GetAllAreas() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 从数据库查询所有区域
	var areas []string
	GVA_DB.Model(&clusterModel.K8sCluster{}).Distinct("area").Where("area != ''").Pluck("area", &areas)
	return areas
}

// Shutdown 关闭所有集群连接
func (m *K8sClusterManager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for clusterId, client := range m.clusters {
		close(client.StopCh)
		if client.PodGroupInformer != nil {
			client.PodGroupInformer.Stop()
		}

		logx.Info("关闭集群连接", clusterId)
	}
	m.clusters = make(map[uint]*ClusterClientInfo)
}

// createClusterClient 根据kubeconfig创建集群客户端
func createClusterClient(clusterId uint, kubeConfig []byte, clusterName, area string, orderManager podgroup.OrderManager) (*ClusterClientInfo, error) {
	// 1. 创建rest.Config
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("解析kubeconfig失败: %w", err)
	}

	// 2. 创建标准客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建clientset失败: %w", err)
	}

	// 3. 创建动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建dynamic client失败: %w", err)
	}

	// 4. 创建controller-runtime client
	scheme := runtime.NewScheme()
	nbv1.AddToScheme(scheme)
	tbv1alpha1.AddToScheme(scheme)
	piperv1beta1.AddToScheme(scheme)

	runtimeClient, err := ctrlclient.New(config, ctrlclient.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("创建runtime client失败: %w", err)
	}

	// 5. 创建Apisix client
	apisixClient, err := apisixclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建apisix client失败: %w", err)
	}

	// 6. 创建Informer工厂和Listers
	factory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	podLister := factory.Core().V1().Pods().Lister()
	stsLister := factory.Apps().V1().StatefulSets().Lister()
	deployLister := factory.Apps().V1().Deployments().Lister()

	// 7. 启动Informer
	stopCh := make(chan struct{})
	factory.Start(stopCh)

	// 等待缓存同步
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	factory.WaitForCacheSync(ctx.Done())

	// 8. 创建Kubeflow CRD Clients
	notebookClient := NewNotebookClient(runtimeClient)
	tensorboardClient := NewTensorboardClient(runtimeClient)

	// 9. 创建Volcano Client
	volcanoClient, err := NewVolcanoClient(config)
	if err != nil {
		return nil, err
	}

	// 10. 创建并启动 PodGroup Informer
	volcanoInformerFactory := vcinformers.NewSharedInformerFactory(volcanoClient.Clientset(), 30*time.Second)
	podGroupInformer := volcanoInformerFactory.Scheduling().V1beta1().PodGroups()
	pgInformer := podgroup.NewPodGroupInformerFactory(
		podGroupInformer,
		clusterId,
		GVA_DB,
		orderManager,
	)
	volcanoInformerFactory.Start(stopCh)
	if err := pgInformer.Start(context.Background()); err != nil {
		logx.Error("启动PodGroup Controller失败", clusterId, err)
	}

	return &ClusterClientInfo{
		ClusterId:         clusterId,
		ClusterName:       clusterName,
		Area:              area,
		ClientSet:         clientset,
		DynamicClient:     dynamicClient,
		RestConfig:        config,
		RuntimeClient:     runtimeClient,
		K8sInformer:       factory,
		PodLister:         podLister,
		StsLister:         stsLister,
		DeployLister:      deployLister,
		NotebookClient:    notebookClient,
		TensorboardClient: tensorboardClient,
		ApisixClient:      apisixClient,
		VolcanoClient:     volcanoClient,
		PodGroupInformer:  pgInformer,
		StopCh:            stopCh,
	}, nil
}
