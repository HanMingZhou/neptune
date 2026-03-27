package initialize

import (
	"gin-vue-admin/global"
	clusterModel "gin-vue-admin/model/cluster"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

// InitK8sClusterManager 初始化多集群管理器（从数据库加载集群配置）
func InitK8sClusterManager() *global.K8sClusterManager {
	manager := global.NewK8sClusterManager()

	// 从数据库加载所有启用的集群配置
	var clusters []clusterModel.K8sCluster
	if err := global.GVA_DB.Where("status = ?", clusterModel.ClusterStatusEnabled).Find(&clusters).Error; err != nil {
		logx.Error("加载集群配置失败", err)
		panic(err)
	}

	// 逐个初始化集群
	for _, cluster := range clusters {
		if err := manager.AddCluster(cluster.ID, []byte(cluster.KubeConfig), cluster.Name, cluster.Area); err != nil {
			logx.Error("初始化集群失败",
				cluster.ID,
				cluster.Name,
				err,
			)
			panic(err)
		}
	}
	return manager
}

// ReloadCluster 重新加载指定集群（用于配置更新后）
func ReloadCluster(manager *global.K8sClusterManager, clusterId uint) error {
	var cluster clusterModel.K8sCluster
	if err := global.GVA_DB.Where("id = ?", clusterId).First(&cluster).Error; err != nil {
		return errors.Wrap(err, "查询集群配置失败")
	}

	// 先移除旧的
	_ = manager.RemoveCluster(clusterId)

	// 如果集群已停用，直接返回
	if cluster.Status != clusterModel.ClusterStatusEnabled {
		return nil
	}

	// 添加新的
	if err := manager.AddCluster(cluster.ID, []byte(cluster.KubeConfig), cluster.Name, cluster.Area); err != nil {
		return errors.Wrap(err, "重新加载集群失败")
	}

	return nil
}
