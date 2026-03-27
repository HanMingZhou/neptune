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
type GetPipeListResp struct {
	PageInfo
	Name       string `json:"name" form:"name"`
	Type       string `json:"type" form:"type"`
	Status     string `json:"status" form:"status"`
	NotebookId int64  `json:"notebookId" form:"notebookId"`
	CreatedBy  string `json:"createdBy" form:"createdBy"`
}

type AddPipeResp struct {
	InstanceName string `json:"instanceName" form:"instanceName" binding:"required"`
	Namespace    string `json:"namespace" form:"namespace"` // 默认为 default
}

type UpdatePipeResp struct {
	Id           uint   `json:"id" form:"id"`
	InstanceName string `json:"instanceName" form:"instanceName" binding:"required"`
	Namespace    string `json:"namespace" form:"namespace"`
}

type DeletePipeResp struct {
	Id           uint   `json:"id" form:"id"`
	InstanceName string `json:"instanceName" form:"instanceName" binding:"required"`
	Namespace    string `json:"namespace" form:"namespace"`
}
