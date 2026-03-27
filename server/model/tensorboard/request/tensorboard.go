package request

import "gorm.io/gorm"

// PageInfo 分页参数
type PageInfo struct {
	Page     int  `json:"page" form:"page"`
	PageSize int  `json:"pageSize" form:"pageSize"`
	Id       uint `json:"id" form:"id"`
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

// GetTensorBoardListReq 获取TensorBoard列表请求
type GetTensorBoardListReq struct {
	PageInfo
}

// AddTensorBoardReq 创建TensorBoard请求
type AddTensorBoardReq struct {
	InstanceName      string            `json:"instanceName" form:"instanceName" binding:"required"`           // 实例名称
	Namespace         string            `json:"namespace" form:"namespace" binding:"required"`                 // 所属命名空间
	OwnerName         string            `json:"ownerName" form:"ownerName" binding:"required"`                 // 所属实例名称
	OwnerType         string            `json:"ownerType" form:"ownerType" binding:"required"`                 // 所属实例类型
	LogsPath          string            `json:"logsPath" form:"logsPath" binding:"required"`                   // 日志路径
	EnableTensorboard bool              `json:"enableTensorboard" form:"enableTensorboard" binding:"required"` // 是否启用 Tensorboard
	Labels            map[string]string `json:"labels" form:"labels" binding:"required"`                       // 标签
}

// UpdateTensorBoardReq 更新TensorBoard请求
type UpdateTensorBoardReq struct {
	Id                uint              `json:"id" form:"id"`
	InstanceName      string            `json:"instanceName" form:"instanceName" binding:"required"`           // 实例名称
	Namespace         string            `json:"namespace" form:"namespace" binding:"required"`                 // 所属命名空间
	OwnerName         string            `json:"ownerName" form:"ownerName" binding:"required"`                 // 所属实例名称
	OwnerType         string            `json:"ownerType" form:"ownerType" binding:"required"`                 // 所属实例类型
	LogsPath          string            `json:"logsPath" form:"logsPath"`                                      // 日志路径
	EnableTensorboard bool              `json:"enableTensorboard" form:"enableTensorboard" binding:"required"` // 是否启用 Tensorboard
	Labels            map[string]string `json:"labels" form:"labels" binding:"required"`                       // 标签
}

// DeleteTensorBoardReq 删除TensorBoard请求
type DeleteTensorBoardReq struct {
	Id           uint   `json:"id" form:"id"`
	InstanceName string `json:"instanceName" form:"instanceName" binding:"required"` // 实例名称
	OwnerName    string `json:"ownerName" form:"ownerName" binding:"required"`       // 所属实例名称
	OwnerType    string `json:"ownerType" form:"ownerType" binding:"required"`       // 所属实例类型
	Namespace    string `json:"namespace" form:"namespace"`                          // 所属命名空间
}
