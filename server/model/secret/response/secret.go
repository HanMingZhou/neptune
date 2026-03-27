package response

import "gorm.io/gorm"

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int  `json:"page" form:"page"`         // 页码,默认1
	PageSize int  `json:"pageSize" form:"pageSize"` // 每页大小,默认10
	Id       uint `json:"id" form:"id"`             // 关键字
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

// 获取分类列表
type GetPVCListReq struct {
	PageInfo
	Name       string `json:"name" form:"name"`
	Type       string `json:"type" form:"type"`
	Status     string `json:"status" form:"status"`
	NotebookId uint   `json:"notebookId" form:"notebookId"`
	CreatedBy  string `json:"createdBy" form:"createdBy"`
}

type AddPVCReq struct {
	InstanceName string           `json:"instanceName" form:"instanceName" binding:"required"`
	Namespace    string           `json:"namespace" form:"namespace"`       // 默认为 default
	VolumeMounts []VolumeMountReq `json:"volumeMounts" form:"volumeMounts"` // 挂载列表
}
type VolumeMountReq struct {
	Name       string `json:"name"`       // 卷名称
	MountsPath string `json:"mountsPath"` // 容器内路径
	Size       string `json:"size"`       // 大小 (如 10Gi)
	Type       string `json:"type"`       // 类型 (如 dataset, model, workspace)
}

type UpdatePVCReq struct {
	Id           uint   `json:"id" form:"id"`
	InstanceName string `json:"name" form:"name" binding:"required"`
	Image        string `json:"image" form:"image"`
	CPU          string `json:"cpu" form:"cpu"`
	Memory       string `json:"memory" form:"memory"`
}

type DeletePVCReq struct {
	Id           uint   `json:"id" form:"id"`
	InstanceName string `json:"name" form:"name" binding:"required"`
	Namespace    string `json:"namespace" form:"namespace"`
}
