package request

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
type GetSecretListReq struct {
	PageInfo
	Name      string `json:"name" form:"name"`           // 实例名称
	Type      string `json:"type" form:"type"`           // 实例类型
	Status    string `json:"status" form:"status"`       // 状态
	CreatedBy string `json:"createdBy" form:"createdBy"` // 创建人
}

type AddSecretReq struct {
	InstanceName string            `json:"instanceName" form:"instanceName" binding:"required"` // 实例名称
	Namespace    string            `json:"namespace" form:"namespace"`                          // 默认为 default
	Content      map[string]string `json:"content" form:"content"`                              // 内容
	InstanceType string            `json:"instanceType" form:"instanceType"`                    // 实例类型
}

type UpdateSecretReq struct {
	Id           uint              `json:"id" form:"id"`
	InstanceName string            `json:"name" form:"name" binding:"required"` // 实例名称
	Content      map[string]string `json:"content" form:"content"`              // 内容
	InstanceType string            `json:"instanceType" form:"instanceType"`    // 实例类型
}

type DeleteSecretReq struct {
	Id           uint              `json:"id" form:"id"`
	InstanceName string            `json:"name" form:"name" binding:"required"` // 实例名称
	Namespace    string            `json:"namespace" form:"namespace"`          // 命名空间
	Content      map[string]string `json:"content" form:"content"`              // 内容
	InstanceType string            `json:"instanceType" form:"instanceType"`    // 实例类型
}
