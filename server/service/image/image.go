package image

import (
	"context"
	"errors"
	"gin-vue-admin/global"
	imageModel "gin-vue-admin/model/image"
	"gin-vue-admin/model/image/request"
	"gin-vue-admin/model/image/response"
	sysModel "gin-vue-admin/model/system"
	"gin-vue-admin/utils/uuid"
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

	// 转换为响应格式
	list := make([]response.ImageItem, 0, len(images))
	for _, img := range images {
		list = append(list, response.ImageItem{
			ID:         img.ID,
			Name:       img.Name,
			Image:      img.ImageAddr,
			Size:       img.Size,
			ImageUUID:  img.ImageUUID,
			Area:       img.Area,
			CreateTime: img.CreatedAt.Format("2006-01-02 15:04:05"),
			UsageType:  img.UsageType,
			Type:       img.Type,
			ImagePath:  img.ImagePath,
			UserId:     img.UserId,
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
	img := imageModel.Image{
		Name:      req.Name,
		UserId:    userId,
		Area:      req.Area,
		UsageType: req.UsageType,
		Type:      imgType,
		ImageAddr: req.ImageAddr,
		Size:      req.Size,
		ImagePath: req.ImagePath,
		ImageUUID: uuid.GenerateUUID(),
	}
	return global.GVA_DB.Create(&img).Error
}

// UpdateImage 更新镜像
func (s *ImageService) UpdateImage(ctx context.Context, req *request.UpdateImageReq, userId uint) error {
	// 判断镜像是否为自己创建
	var img imageModel.Image
	if err := global.GVA_DB.Where("id = ?", req.Id).First(&img).Error; err != nil {
		return err
	}
	if img.Type == imageModel.ImageTypeSystem {
		if img.UserId != userId {
			return errors.New("系统镜像不属于自己")
		}
	}
	updates := map[string]interface{}{
		"name":       req.Name,
		"type":       req.Type,
		"usage_type": req.UsageType,
		"image_addr": req.ImageAddr,
		"area":       req.Area,
		"size":       req.Size,
		"image_path": req.ImagePath,
	}
	return global.GVA_DB.Model(&imageModel.Image{}).Where("id = ?", req.Id).Updates(updates).Error
}
