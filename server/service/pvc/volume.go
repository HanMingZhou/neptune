package pvc

import (
	"context"
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model/cluster"
	"gin-vue-admin/model/consts"
	productModel "gin-vue-admin/model/product"
	pvcModel "gin-vue-admin/model/pvc"
	"gin-vue-admin/model/pvc/request"
	"gin-vue-admin/model/pvc/response"
	"gin-vue-admin/utils"
	helper "gin-vue-admin/utils/k8s"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
VolumeService 提供文件存储（PVC）的 CRUD 操作

功能：
1. CreateVolume - 创建 K8s PVC + 数据库记录
2. GetVolumeList - 查询列表（联查使用状态）
3. ExpandVolume - 扩容 PVC
4. DeleteVolume - 删除前校验是否被使用
5. GetAreaList - 获取可用区域列表
*/

type VolumeManager interface {
	CreateVolume(ctx context.Context, req *request.CreateVolumeReq) error
	GetVolumeList(ctx context.Context, req *request.GetVolumeListReq) (*response.VolumeListResp, error)
	ExpandVolume(ctx context.Context, req *request.ExpandVolumeReq) error
	DeleteVolume(ctx context.Context, req *request.DeleteVolumeReq) error
	GetAreaList(ctx context.Context) (*response.AreaListResp, error)
}

var _ VolumeManager = (*VolumeService)(nil)

type VolumeService struct{}

var VolumeServiceApp = new(VolumeService)

// CreateVolume 创建文件存储
func (v *VolumeService) CreateVolume(ctx context.Context, req *request.CreateVolumeReq) error {
	// 1. 验证参数
	if req.Name == "" {
		return errors.New("存储名称不能为空")
	}
	if req.Size <= 0 {
		return errors.New("容量必须大于0")
	}
	if req.Area == "" {
		return errors.New("区域不能为空")
	}

	// 2. 获取产品和集群
	if req.ProductId <= 0 {
		return errors.New("存储产品ID不能为空")
	}
	var prod productModel.Product
	if err := global.GVA_DB.Where("id = ? AND product_type = ?", req.ProductId, productModel.ProductTypeStorage).First(&prod).Error; err != nil {
		logx.Error("获取存储产品失败", err)
		return errors.New("指定的存储产品不存在")
	}
	if prod.Status != productModel.ProductStatusEnabled {
		return errors.New("指定的存储产品已下架")
	}
	if req.ClusterId <= 0 {
		return errors.New("集群ID不能为空")
	}
	if prod.ClusterId != req.ClusterId {
		return errors.New("存储产品所属集群与所选集群不匹配")
	}

	var cluster cluster.K8sCluster
	err := global.GVA_DB.Where("id = ?", req.ClusterId).First(&cluster).Error
	if err != nil {
		logx.Error("获取集群失败", err)
		return errors.New("指定的集群不可用")
	}
	clientInfo := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(req.ClusterId)
	if clientInfo == nil {
		logx.Error("获取集群失败", logx.Field("clusterId", req.ClusterId))
		return errors.New("指定的集群不可用")
	}
	clientSet := clientInfo.ClientSet

	// 3. 系统生成 PVC 名称 (vol-时间戳-随机数)，确保符合 K8s 命名规范且唯一（全部小写）
	pvcName := strings.ToLower(fmt.Sprintf("vol-%d-%s", time.Now().Unix(), utils.RandomString(4)))

	// 4. 确保 namespace 存在
	if err := helper.EnsureNamespace(ctx, clientSet, req.Namespace); err != nil {
		return err
	}

	// 5. 创建 K8s PVC
	storageClassName := prod.StorageClass
	if storageClassName == "" {
		storageClassName = cluster.StorageClass
	}
	if storageClassName == "" {
		storageClassName = consts.DefaultStorageClass // 默认 StorageClass
	}

	sizeStr := fmt.Sprintf("%dGi", req.Size)
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"app": "neptune-volume",
			},
			Annotations: map[string]string{
				"area":       req.Area,
				"cluster-id": cluster.Name,
				"product-id": fmt.Sprintf("%d", prod.ID),
			},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany}, // RWX
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(sizeStr),
				},
			},
			StorageClassName: &storageClassName,
		},
	}

	if _, err := clientSet.CoreV1().PersistentVolumeClaims(req.Namespace).Create(ctx, pvc, metav1.CreateOptions{}); err != nil {
		logx.Error("创建K8s PVC失败", err)
		return errors.Wrap(err, "创建存储失败")
	}

	// 6. 保存到数据库
	volume := &pvcModel.Volume{
		Name:      req.Name,
		Namespace: req.Namespace,
		Size:      req.Size,
		Type:      req.Type,
		Status:    consts.PVCStatusReady,
		PVCName:   pvcName,
		ClusterId: cluster.ID,
		ProductId: prod.ID,
		UserId:    req.UserId,
	}

	if err := global.GVA_DB.Create(volume).Error; err != nil {
		// 回滚：删除已创建的 PVC
		_ = clientSet.CoreV1().PersistentVolumeClaims(req.Namespace).Delete(ctx, pvcName, metav1.DeleteOptions{})
		logx.Error("保存Volume到数据库失败", err)
		return errors.Wrap(err, "保存存储信息失败")
	}

	logx.Info("创建Volume成功", logx.Field("name", req.Name), logx.Field("pvcName", pvcName))
	return nil
}

// GetVolumeList 获取文件存储列表
func (v *VolumeService) GetVolumeList(ctx context.Context, req *request.GetVolumeListReq) (*response.VolumeListResp, error) {
	var volumes []pvcModel.Volume
	var total int64

	db := global.GVA_DB.Model(&pvcModel.Volume{})

	// 根据筛选条件查询
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	if req.Area != "" {
		db = db.Where("area = ?", req.Area)
	}
	if req.UserId > 0 {
		db = db.Where("user_id = ?", req.UserId)
	}
	if req.ClusterId > 0 {
		db = db.Where("cluster_id = ?", req.ClusterId)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	if err := db.Scopes(req.Paginate()).Preload("K8sCluster").Preload("Product").Order("created_at DESC").Find(&volumes).Error; err != nil {
		return nil, err
	}

	// 构建响应
	list := make([]response.VolumeItem, 0, len(volumes))
	for _, vol := range volumes {
		productName := ""
		if vol.Product != nil {
			productName = vol.Product.Name
		}
		item := response.VolumeItem{
			ID:          vol.ID,
			Name:        vol.Name,
			Size:        vol.Size,
			Type:        vol.Type,
			Status:      v.getVolumeStatus(ctx, &vol),
			CreatedAt:   vol.CreatedAt.Format("2006-01-02 15:04:05"),
			UsedBy:      v.getVolumeUsage(&vol),
			ClusterId:   vol.K8sCluster.ID,
			ProductId:   vol.ProductId,
			ProductName: productName,
			Area:        vol.K8sCluster.Area,
		}
		list = append(list, item)
	}

	return &response.VolumeListResp{
		List:  list,
		Total: total,
	}, nil
}

// ExpandVolume 扩容文件存储
func (v *VolumeService) ExpandVolume(ctx context.Context, req *request.ExpandVolumeReq) error {
	// 1. 查询 Volume
	var volume pvcModel.Volume
	if err := global.GVA_DB.Where("id = ? AND user_id = ?", req.Id, req.UserId).First(&volume).Error; err != nil {
		return errors.New("存储不存在或无权操作")
	}

	// 2. 获取当前容量
	currentSize := volume.Size

	// 3. 校验：只能扩大
	if req.Size <= currentSize {
		return errors.New("新容量必须大于当前容量")
	}

	// 4. 获取集群客户端 (优先使用数据库记录的 ClusterId)
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(volume.ClusterId)
	if cluster == nil {
		return errors.New("对应集群不可用")
	}
	clientSet := cluster.ClientSet

	// 5. 更新 K8s PVC
	newSizeStr := fmt.Sprintf("%dGi", req.Size)
	pvc, err := clientSet.CoreV1().PersistentVolumeClaims(volume.Namespace).Get(ctx, volume.PVCName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "获取PVC失败")
	}

	pvc.Spec.Resources.Requests[corev1.ResourceStorage] = resource.MustParse(newSizeStr)
	if _, err := clientSet.CoreV1().PersistentVolumeClaims(volume.Namespace).Update(ctx, pvc, metav1.UpdateOptions{}); err != nil {
		logx.Error("扩容K8s PVC失败", err)
		return errors.Wrap(err, "扩容失败")
	}

	// 6. 更新数据库
	if err := global.GVA_DB.Model(&pvcModel.Volume{}).Where("id = ?", req.Id).Update("size", req.Size).Error; err != nil {
		return errors.Wrap(err, "更新存储信息失败")
	}

	logx.Info("扩容Volume成功", logx.Field("id", req.Id), logx.Field("newSize", req.Size))
	return nil
}

// DeleteVolume 删除文件存储
func (v *VolumeService) DeleteVolume(ctx context.Context, req *request.DeleteVolumeReq) error {
	// 1. 查询 Volume
	var volume pvcModel.Volume
	if err := global.GVA_DB.Where("id = ? AND user_id = ?", req.Id, req.UserId).First(&volume).Error; err != nil {
		return errors.New("存储不存在或无权操作")
	}

	// 2. 检查是否被使用
	usedBy := v.getVolumeUsage(&volume)
	if len(usedBy) > 0 {
		names := ""
		for _, u := range usedBy {
			names += fmt.Sprintf("%s: %s, ", u.Type, u.Name)
		}
		return errors.New("该存储正在被使用: " + names + "请先停止相关任务")
	}

	// 3. 获取集群客户端 (使用数据库记录的 ClusterId)
	cluster := global.GVA_K8S_CLUSTER_MANAGER.GetCluster(volume.ClusterId)
	if cluster == nil {
		// 如果集群已经不存在了，直接删除数据库记录
		logx.Errorf("集群不存在，直接删除数据库记录", logx.Field("clusterId", volume.ClusterId))
	} else {
		clientSet := cluster.ClientSet

		// 4. 删除 K8s PVC
		if err := clientSet.CoreV1().PersistentVolumeClaims(volume.Namespace).Delete(ctx, volume.PVCName, metav1.DeleteOptions{}); err != nil {
			logx.Error("删除K8s PVC失败", err)
			// 继续删除数据库记录
		}
	}

	// 5. 删除数据库记录 (恢复为软删除)
	if err := global.GVA_DB.Delete(&pvcModel.Volume{}, req.Id).Error; err != nil {
		return errors.Wrap(err, "删除存储记录失败")
	}

	logx.Info("删除Volume成功", logx.Field("id", req.Id))
	return nil
}

// GetAreaList 获取可用区域/集群列表
func (v *VolumeService) GetAreaList(ctx context.Context) (*response.AreaListResp, error) {
	clusters := global.GVA_K8S_CLUSTER_MANAGER.ListClusters()
	list := make([]response.ClusterInfo, 0, len(clusters))
	for _, c := range clusters {
		list = append(list, response.ClusterInfo{
			ID:   uint(c.ClusterId),
			Name: c.ClusterName,
			Area: c.Area,
		})
	}
	return &response.AreaListResp{
		Clusters: list,
	}, nil
}

// ============ 辅助函数 ============

// getVolumeStatus 获取 Volume 使用状态
func (v *VolumeService) getVolumeStatus(ctx context.Context, vol *pvcModel.Volume) string {
	usedBy := v.getVolumeUsage(vol)
	if len(usedBy) > 0 {
		return "使用中"
	}
	return "未使用"
}

// getVolumeUsage 查询 Volume 被哪些资源使用
func (v *VolumeService) getVolumeUsage(vol *pvcModel.Volume) []response.VolumeUsage {
	var usages []response.VolumeUsage

	// 1. 查询 Notebooks
	var notebookNames []string
	global.GVA_DB.Table("notebook_volumes").
		Joins("JOIN notebooks ON notebooks.id = notebook_volumes.notebook_id").
		Where("notebook_volumes.pvc_id = ? AND notebook_volumes.deleted_at IS NULL AND notebooks.deleted_at IS NULL", vol.ID).
		Pluck("COALESCE(notebooks.display_name, notebooks.instance_name)", &notebookNames)

	for _, name := range notebookNames {
		usages = append(usages, response.VolumeUsage{
			Type: "Notebook",
			Name: name,
		})
	}

	// 2. 查询 TrainingJobs
	var trainingDisplayNames []string
	global.GVA_DB.Table("training_job_mounts").
		Joins("JOIN training_jobs ON training_jobs.id = training_job_mounts.job_id").
		Where("training_job_mounts.pvc_id = ? AND training_job_mounts.deleted_at IS NULL AND training_jobs.deleted_at IS NULL", vol.ID).
		Pluck("COALESCE(training_jobs.display_name, training_jobs.job_name)", &trainingDisplayNames)

	for _, name := range trainingDisplayNames {
		usages = append(usages, response.VolumeUsage{
			Type: "TrainingJob",
			Name: name,
		})
	}

	return usages
}
