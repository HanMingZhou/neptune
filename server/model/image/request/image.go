package request

import (
	"gorm.io/gorm"
)

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int   `json:"page" form:"page"`         // 页码,默认1
	PageSize int   `json:"pageSize" form:"pageSize"` // 每页大小,默认10
	Id       int64 `json:"id" form:"id"`             // 关键字
}

func (r *PageInfo) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.Page <= 0 {
			r.Page = 1
		}
		switch {
		case r.PageSize > 100:
			r.PageSize = 100
		case r.PageSize <= 0:
			r.PageSize = 10
		}
		offset := (r.Page - 1) * r.PageSize
		return db.Offset(offset).Limit(r.PageSize)
	}
}

// 获取镜像列表
type GetImageListReq struct {
	PageInfo
	Name       string `json:"name" form:"name"`             // 镜像名称
	Type       int64  `json:"type" form:"type"`             // 1 系统镜像，2 自定义镜像
	UsageType  int64  `json:"usageType" form:"usageType"`   // 1 容器实例，2 训练镜像，3 推理镜像
	ImageId    uint   `json:"imageId" form:"imageId"`       // 镜像ID
	Size       string `json:"size" form:"size"`             // 镜像大小
	ImageUUID  string `json:"imageUUID" form:"imageUUID"`   // 镜像UUID
	Area       string `json:"area" form:"area"`             // 镜像区域
	CreateTime string `json:"createTime" form:"createTime"` // 创建时间
}

// 添加自定义镜像
type AddImageReq struct {
	Name      string `json:"name" form:"name" binding:"required"`           // 镜像名称
	Area      string `json:"area" form:"area" binding:"required"`           // 镜像区域
	UsageType int64  `json:"usageType" form:"usageType" binding:"required"` // 1 容器实例，2 训练镜像，3 推理镜像
	Type      int64  `json:"type" form:"type"`                              // 1 系统镜像，2 自定义镜像
	ImageAddr string `json:"imageAddr" form:"imageAddr"`                    // 镜像地址
	Size      string `json:"size" form:"size"`                              // 镜像大小
	ImagePath string `json:"imagePath" form:"imagePath"`                    // 镜像路径
	ImageUUID string `json:"imageUUID" form:"imageUUID"`                    // 镜像UUID
}

// 更新镜像
type UpdateImageReq struct {
	Id        uint   `json:"id" form:"id" binding:"required"`
	Name      string `json:"name" form:"name"`
	Type      int64  `json:"type" form:"type"`
	UsageType int64  `json:"usageType" form:"usageType"`
	ImageAddr string `json:"imageAddr" form:"imageAddr"`
	Area      string `json:"area" form:"area"`
	Size      string `json:"size" form:"size"`
	ImagePath string `json:"imagePath" form:"imagePath"`
}

// 删除镜像
type DeleteImageReq struct {
	Id uint `json:"id" form:"id"`
}
