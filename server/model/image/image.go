package image

import (
	"gorm.io/gorm"
)

// Image 类型常量
const (
	ImageTypeSystem = 1 // 系统镜像
	ImageTypeCustom = 2 // 自定义镜像

	ImageUsageTypeNotebook = 1 // 容器实例
	ImageUsageTypeTrain    = 2 // 训练镜像
	ImageUsageTypeInfer    = 3 // 推理镜像
)

// Image 数据库模型，用于存储Image的元数据
type Image struct {
	gorm.Model
	Name      string `json:"name" gorm:"column:name;comment:镜像名称"`
	Type      int64  `json:"type" gorm:"column:type;comment:镜像类型:1系统镜像2自定义镜像;default:1"`
	UsageType int64  `json:"usageType" gorm:"column:usage_type;comment:镜像用途:1容器实例2训练镜像3推理镜像"`
	UserId    uint   `json:"userId" gorm:"column:user_id;comment:用户ID"`
	ImagePath string `json:"imagePath" gorm:"column:image_path;comment:镜像路径"`
	ImageAddr string `json:"imageAddr" gorm:"column:image_addr;comment:镜像地址"`
	Size      string `json:"size" gorm:"column:size;comment:镜像大小"`
	ImageUUID string `json:"imageUUID" gorm:"column:image_uuid;comment:镜像UUID"`
	Area      string `json:"area" gorm:"column:area;comment:镜像区域"`
}

// TableName 指定表名
func (Image) TableName() string {
	return "images"
}
