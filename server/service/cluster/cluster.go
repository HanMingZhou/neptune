package cluster

import (
	"context"
	"gin-vue-admin/global"
	clusterModel "gin-vue-admin/model/cluster"
	"gin-vue-admin/model/cluster/request"
	"gin-vue-admin/model/cluster/response"
	"gin-vue-admin/model/consts"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterService struct{}

// GetClusterList 获取集群列表
func (s *ClusterService) GetClusterList(req *request.GetClusterListReq) (*response.ClusterListResponse, error) {
	var clusters []clusterModel.K8sCluster
	var total int64

	db := global.GVA_DB.Model(&clusterModel.K8sCluster{})

	// 关键词搜索
	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR area LIKE ? OR description LIKE ? OR api_server LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 状态筛选
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, errors.Wrap(err, "查询集群总数失败")
	}

	if err := db.Order("id DESC").Find(&clusters).Error; err != nil {
		return nil, errors.Wrap(err, "查询集群列表失败")
	}

	list := make([]response.ClusterItem, 0, len(clusters))
	for _, c := range clusters {
		// 查询集群节点总数和Master内网IP
		nodeCount := 0
		internalIP := ""
		clusterClient := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(c.ID)
		if clusterClient != nil && clusterClient.ClientSet != nil {
			nodeList, err := clusterClient.ClientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
			if err == nil {
				nodeCount = len(nodeList.Items)
				// 获取第一个 Master 节点的内网IP
				for _, node := range nodeList.Items {
					if _, ok := node.Labels[consts.LabelNodeRoleMaster]; !ok {
						if _, ok2 := node.Labels[consts.LabelNodeRoleControlPlane]; !ok2 {
							continue
						}
					}
					for _, addr := range node.Status.Addresses {
						if addr.Type == corev1.NodeInternalIP {
							internalIP = addr.Address
							break
						}
					}
					if internalIP != "" {
						break
					}
				}
			}
		}

		list = append(list, response.ClusterItem{
			ID:           c.ID,
			Name:         c.Name,
			Area:         c.Area,
			Description:  c.Description,
			KubeConfig:   c.KubeConfig,
			ApiServer:    c.ApiServer,
			Status:       c.Status,
			HarborAddr:   c.HarborAddr,
			StorageClass: c.StorageClass,
			NodeCount:    nodeCount,
			InternalIP:   internalIP,
			CreatedAt:    c.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &response.ClusterListResponse{
		List:  list,
		Total: total,
	}, nil
}

// CreateCluster 创建集群
func (s *ClusterService) CreateCluster(req *request.CreateClusterReq) error {
	// 检查名称是否重复
	var existing clusterModel.K8sCluster
	if err := global.GVA_DB.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		return errors.New("集群名称已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.Wrap(err, "查询集群失败")
	}

	cluster := clusterModel.K8sCluster{
		Name:         req.Name,
		Area:         req.Area,
		Description:  req.Description,
		KubeConfig:   req.KubeConfig,
		ApiServer:    req.ApiServer,
		Status:       req.Status,
		HarborAddr:   req.HarborAddr,
		StorageClass: req.StorageClass,
	}

	if err := global.GVA_DB.Create(&cluster).Error; err != nil {
		return errors.Wrap(err, "创建集群失败")
	}

	return nil
}

// UpdateCluster 更新集群
func (s *ClusterService) UpdateCluster(req *request.UpdateClusterReq) error {
	var existing clusterModel.K8sCluster
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&existing).Error; err != nil {
		return errors.Wrap(err, "集群不存在")
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		// 检查名称是否与其他集群重复
		var other clusterModel.K8sCluster
		err := global.GVA_DB.Where("name = ? AND id != ?", req.Name, req.ID).First(&other).Error
		if err == nil {
			return errors.New("集群名称已存在")
		}
		updates["name"] = req.Name
	}
	if req.Area != "" {
		updates["area"] = req.Area
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.KubeConfig != "" {
		updates["kubeconfig"] = req.KubeConfig
	}
	if req.ApiServer != "" {
		updates["api_server"] = req.ApiServer
	}
	updates["status"] = req.Status
	if req.HarborAddr != "" {
		updates["harbor_addr"] = req.HarborAddr
	}
	if req.StorageClass != "" {
		updates["storage_class"] = req.StorageClass
	}

	if err := global.GVA_DB.Model(&clusterModel.K8sCluster{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		return errors.Wrap(err, "更新集群失败")
	}

	// 关键字段变更时重载集群连接
	needReload := (req.KubeConfig != "" && req.KubeConfig != existing.KubeConfig) ||
		(req.ApiServer != "" && req.ApiServer != existing.ApiServer) ||
		(req.StorageClass != "" && req.StorageClass != existing.StorageClass) ||
		(req.Area != "" && req.Area != existing.Area)

	if needReload {
		// 取更新后的值
		name := existing.Name
		if req.Name != "" {
			name = req.Name
		}
		area := existing.Area
		if req.Area != "" {
			area = req.Area
		}
		kubeConfig := existing.KubeConfig
		if req.KubeConfig != "" {
			kubeConfig = req.KubeConfig
		}

		if err := global.GVA_K8S_CLUSTER_MANAGER.ReloadCluster(req.ID, []byte(kubeConfig), name, area); err != nil {
			return errors.Wrap(err, "重载集群连接失败")
		}

		// 更新 StorageClass 到内存中的集群信息
		if req.StorageClass != "" {
			if client := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(req.ID); client != nil {
				client.DefaultStorageClass = req.StorageClass
			}
		}
	}

	return nil
}

// DeleteCluster 删除集群
func (s *ClusterService) DeleteCluster(req *request.DeleteClusterReq) error {
	if err := global.GVA_DB.Delete(&clusterModel.K8sCluster{}, req.ID).Error; err != nil {
		return errors.Wrap(err, "删除集群失败")
	}
	return nil
}
