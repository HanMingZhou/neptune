package image

import (
	"context"
	"errors"
	"gin-vue-admin/global"
	clusterModel "gin-vue-admin/model/cluster"
	imageModel "gin-vue-admin/model/image"
	"gin-vue-admin/model/image/request"
	"gin-vue-admin/model/image/response"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/utils/uuid"
	"strings"
)

type ImageService struct{}

// GetImageList 获取镜像列表
func (s *ImageService) GetImageList(ctx context.Context, req *request.GetImageListReq, userId uint) (*response.GetImageListResp, error) {
	var images []imageModel.Image
	var total int64

	db := global.GVA_DB.Model(&imageModel.Image{})

	// 条件筛选
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Type > 0 {
		db = db.Where("type = ?", req.Type)
		if req.Type == imageModel.ImageTypeCustom {
			db = db.Where("user_id = ?", userId)
		}
	}
	if req.ImageId > 0 {
		db = db.Where("id = ?", req.ImageId)
	}
	if req.ClusterId > 0 {
		db = db.Where("cluster_id = ?", req.ClusterId)
	}
	if req.Area != "" {
		db = db.Where("area = ?", req.Area)
	}
	if req.ImageUUID != "" {
		db = db.Where("image_uuid = ?", req.ImageUUID)
	}
	if req.UsageType > 0 {
		db = db.Where("usage_type = ?", req.UsageType)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	if err := db.Scopes(req.Paginate()).Order("created_at DESC").Find(&images).Error; err != nil {
		return nil, err
	}

	clusterNames, err := s.loadClusterNames(ctx, images)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]response.ImageItem, 0, len(images))
	for _, img := range images {
		list = append(list, response.ImageItem{
			ID:          img.ID,
			Name:        img.Name,
			Image:       img.ImageAddr,
			Size:        img.Size,
			ImageUUID:   img.ImageUUID,
			Area:        img.Area,
			ClusterId:   img.ClusterId,
			ClusterName: clusterNames[img.ClusterId],
			CreateTime:  img.CreatedAt.Format("2006-01-02 15:04:05"),
			UsageType:   img.UsageType,
			Type:        img.Type,
			ImagePath:   img.ImagePath,
			UserId:      img.UserId,
		})
	}

	return &response.GetImageListResp{
		List:  list,
		Total: total,
	}, nil
}

// DeleteImage 删除镜像
func (s *ImageService) DeleteImage(ctx context.Context, req *request.DeleteImageReq, userId uint) error {
	if req.Id <= 0 {
		return nil
	}
	// 判断镜像是否为自己创建
	var img imageModel.Image
	if err := global.GVA_DB.Where("id = ?", req.Id).First(&img).Error; err != nil {
		return err
	}
	if img.UserId != userId {
		return errors.New("镜像不属于自己")
	}
	// 从数据库删除记录
	if err := global.GVA_DB.Delete(&imageModel.Image{}, req.Id).Error; err != nil {
		return err
	}

	return nil
}

// CreateImage 创建镜像
func (s *ImageService) CreateImage(ctx context.Context, req *request.AddImageReq, userId uint) error {
	imgType := req.Type
	if imgType == 0 {
		imgType = imageModel.ImageTypeCustom
	}
	if imgType == imageModel.ImageTypeSystem {
		var user sysModel.SysUser
		if err := global.GVA_DB.Where("id = ?", userId).First(&user).Error; err != nil {
			return err
		}
		if !sysModel.SysAuthorityIds[user.AuthorityId] {
			return errors.New("系统镜像不能创建")
		}
	}

	cluster, err := s.getEnabledCluster(ctx, req.ClusterId)
	if err != nil {
		return err
	}
	imageAddr, err := resolveImageAddr(
		req.ImageAddr,
		cluster.HarborAddr,
		req.ImagePath,
	)
	if err != nil {
		return err
	}

	img := imageModel.Image{
		Name:      req.Name,
		UserId:    userId,
		ClusterId: req.ClusterId,
		Area:      cluster.Area,
		UsageType: req.UsageType,
		Type:      imgType,
		ImageAddr: imageAddr,
		Size:      req.Size,
		ImagePath: strings.Trim(strings.TrimSpace(req.ImagePath), "/"),
		ImageUUID: uuid.GenerateUUID(),
	}
	return global.GVA_DB.WithContext(ctx).Create(&img).Error
}

// UpdateImage 更新镜像
func (s *ImageService) UpdateImage(ctx context.Context, req *request.UpdateImageReq, userId uint) error {
	// 判断镜像是否为自己创建
	var img imageModel.Image
	if err := global.GVA_DB.WithContext(ctx).Where("id = ?", req.Id).First(&img).Error; err != nil {
		return err
	}
	if img.Type == imageModel.ImageTypeSystem {
		if img.UserId != userId {
			return errors.New("系统镜像不属于自己")
		}
	}

	clusterId := img.ClusterId
	if req.ClusterId > 0 {
		clusterId = req.ClusterId
	}
	cluster, err := s.getEnabledCluster(ctx, clusterId)
	if err != nil {
		return err
	}

	name := img.Name
	if strings.TrimSpace(req.Name) != "" {
		name = req.Name
	}

	imageType := img.Type
	if req.Type > 0 {
		imageType = req.Type
	}

	usageType := img.UsageType
	if req.UsageType > 0 {
		usageType = req.UsageType
	}

	imagePath := img.ImagePath
	if strings.TrimSpace(req.ImagePath) != "" {
		imagePath = strings.Trim(strings.TrimSpace(req.ImagePath), "/")
	}

	imageAddr, err := buildImageAddr(cluster.HarborAddr, imagePath)
	if err != nil {
		return err
	}
	if strings.TrimSpace(req.ImageAddr) != "" {
		imageAddr = strings.TrimSpace(req.ImageAddr)
	}

	size := img.Size
	if strings.TrimSpace(req.Size) != "" {
		size = req.Size
	}

	updates := map[string]interface{}{
		"name":       name,
		"cluster_id": clusterId,
		"type":       imageType,
		"usage_type": usageType,
		"image_addr": imageAddr,
		"area":       cluster.Area,
		"size":       size,
		"image_path": imagePath,
	}
	return global.GVA_DB.WithContext(ctx).Model(&imageModel.Image{}).Where("id = ?", req.Id).Updates(updates).Error
}

func (s *ImageService) loadClusterNames(ctx context.Context, images []imageModel.Image) (map[uint]string, error) {
	clusterIDs := make([]uint, 0)
	seen := make(map[uint]struct{})
	for _, img := range images {
		if img.ClusterId == 0 {
			continue
		}
		if _, ok := seen[img.ClusterId]; ok {
			continue
		}
		seen[img.ClusterId] = struct{}{}
		clusterIDs = append(clusterIDs, img.ClusterId)
	}
	if len(clusterIDs) == 0 {
		return map[uint]string{}, nil
	}

	var clusters []clusterModel.K8sCluster
	if err := global.GVA_DB.WithContext(ctx).Where("id IN ?", clusterIDs).Find(&clusters).Error; err != nil {
		return nil, err
	}

	result := make(map[uint]string, len(clusters))
	for _, cluster := range clusters {
		result[cluster.ID] = cluster.Name
	}
	return result, nil
}

func (s *ImageService) getEnabledCluster(ctx context.Context, clusterId uint) (*clusterModel.K8sCluster, error) {
	if clusterId == 0 {
		return nil, errors.New("请选择集群")
	}

	var cluster clusterModel.K8sCluster
	if err := global.GVA_DB.WithContext(ctx).
		Where("id = ? AND status = ?", clusterId, clusterModel.ClusterStatusEnabled).
		First(&cluster).Error; err != nil {
		return nil, errors.New("集群不存在或已停用")
	}
	return &cluster, nil
}

func buildImageAddr(harborAddr, imagePath string) (string, error) {
	path := strings.Trim(strings.TrimSpace(imagePath), "/")
	if path == "" {
		return "", errors.New("镜像路径不能为空")
	}

	registry := strings.TrimSpace(harborAddr)
	registry = strings.TrimRight(registry, "/")
	if registry == "" {
		return "", errors.New("所选集群未配置 Harbor 地址")
	}

	return registry + "/" + path, nil
}

func resolveImageAddr(imageAddr, harborAddr, imagePath string) (string, error) {
	if strings.TrimSpace(imageAddr) != "" {
		return strings.TrimSpace(imageAddr), nil
	}

	return buildImageAddr(harborAddr, imagePath)
}
