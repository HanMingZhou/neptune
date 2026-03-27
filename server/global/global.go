package global

import (
	"fmt"
	"gin-vue-admin/config"
	"gin-vue-admin/utils/timer"
	"sync"

	apisixclient "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
	"github.com/qiniu/qmgo"
	"github.com/redis/go-redis/v9"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1lister "k8s.io/client-go/listers/apps/v1"
	v1podlister "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"volcano.sh/apis/pkg/client/listers/batch/v1alpha1"
	"volcano.sh/apis/pkg/client/listers/scheduling/v1beta1"
)

var (
	GVA_DB                  *gorm.DB
	GVA_DBList              map[string]*gorm.DB
	GVA_REDIS               redis.UniversalClient
	GVA_REDISList           map[string]redis.UniversalClient
	GVA_MONGO               *qmgo.QmgoClient
	GVA_CONFIG              config.Server
	GVA_VP                  *viper.Viper
	GVA_LOG                 *zap.Logger
	GVA_Timer               timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control             = &singleflight.Group{}
	GVA_ROUTERS             gin.RoutesInfo
	GVA_ACTIVE_DBNAME       *string
	GVA_MCP_SERVER          *server.MCPServer
	BlackCache              local_cache.Cache
	lock                    sync.RWMutex
	GVA_K8S_CLIENT_INFO     *K8sClientInfo     // 单集群模式（已废弃，保留用于兼容）
	GVA_K8S_CLUSTER_MANAGER *K8sClusterManager // 多集群管理器
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return GVA_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := GVA_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}

func GetRedis(name string) redis.UniversalClient {
	redis, ok := GVA_REDISList[name]
	if !ok || redis == nil {
		panic(fmt.Sprintf("redis `%s` no init", name))
	}
	return redis
}

type K8sClientInfo struct {
	ClientSet      *kubernetes.Clientset
	DynamicClient  dynamic.Interface
	RestConfig     *rest.Config
	RuntimeClient  ctrlclient.Client // controller-runtime client，支持直接传原生结构体
	K8sInformer    informers.SharedInformerFactory
	PodLister      v1podlister.PodLister
	StsLister      v1lister.StatefulSetLister
	DeployLister   v1lister.DeploymentLister
	VcLister       v1alpha1.JobLister
	PodGroupLister v1beta1.PodGroupLister
	StopCh         chan struct{}
	// Kubeflow CRD Clients
	NotebookClient    *NotebookClient
	TensorboardClient *TensorboardClient
	ApisixClient      apisixclient.Interface
}
