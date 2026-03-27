package request

import "gorm.io/gorm"

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
	InstanceName string            `json:"instanceName" form:"instanceName" binding:"required"`
	Namespace    string            `json:"namespace" form:"namespace"`       // 默认为 default
	ClusterId    uint              `json:"clusterId" form:"clusterId"`       // 集群ID
	VolumeMounts []*VolumeMountReq `json:"volumeMounts" form:"volumeMounts"` // 挂载列表
}
type VolumeMountReq struct {
	Name       string `json:"name"`       // 卷名称
	MountsPath string `json:"mountsPath"` // 容器内路径
	Size       int64  `json:"size"`       // 大小 (GB)
	Type       string `json:"type"`       // 类型 (如 dataset, model, workspace)
	PVCId      uint   `json:"pvcId"`      // 数据库中 Volume 的 ID
	PVCName    string `json:"pvcName"`    // PVC名称
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

// ==================== Volume 文件存储 ====================

// 创建文件存储
type CreateVolumeReq struct {
	Name      string `json:"name" binding:"required"`      // 存储名称
	Size      int64  `json:"size" binding:"required"`      // 容量(GB)
	Type      int    `json:"type"`                         // 类型(1:dataset, 2:model, 3:workspace)
	Area      string `json:"area" binding:"required"`      // 区域
	ClusterId uint   `json:"clusterId" binding:"required"` // 集群ID (必填)
	ProductId uint   `json:"productId" binding:"required"` // 产品ID (必填)
	UserId    uint   `json:"userId"`                       // 用户ID (必填)
	Namespace string `json:"namespace"`                    // namespace (必填)
}

// 扩容
type ExpandVolumeReq struct {
	Id     uint  `json:"id" binding:"required"`   // 存储 ID
	Size   int64 `json:"size" binding:"required"` // 新容量(GB)，必须大于当前容量
	UserId uint  `json:"userId"`
}

// 删除
type DeleteVolumeReq struct {
	Id     uint `json:"id" binding:"required"`
	UserId uint `json:"userId"`
}

// 列表查询
type GetVolumeListReq struct {
	PageInfo
	Name      string `json:"name" form:"name"`
	Status    string `json:"status" form:"status"`
	Area      string `json:"area" form:"area"`
	UserId    uint   `json:"userId" form:"userId"`
	ClusterId uint   `json:"clusterId" form:"clusterId"`
}

type TaskVolumeInfo struct {
	ID        uint
	Name      string
	Size      int64
	UserId    uint
	ClusterId uint
}
