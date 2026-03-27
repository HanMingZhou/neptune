package response

import "gorm.io/gorm"

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page     int  `json:"page" form:"page"`         // 页码,默认1
	PageSize int  `json:"pageSize" form:"pageSize"` // 每页大小,默认10
	ID       uint `json:"id" form:"id"`             // 关键字
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
type GetNotebookListReq struct {
	PageInfo
	Name       string `json:"name" form:"name"`
	Type       string `json:"type" form:"type"`
	Status     string `json:"status" form:"status"`
	NotebookId uint   `json:"notebookId" form:"notebookId"`
	CreatedBy  string `json:"createdBy" form:"createdBy"`
}

type AddNoteBookReq struct {
	Name              string           `json:"name" form:"name" binding:"required"`
	Namespace         string           `json:"namespace" form:"namespace"` // 默认为 default
	Image             string           `json:"image" form:"image"`
	CPU               string           `json:"cpu" form:"cpu"`
	Memory            string           `json:"memory" form:"memory"`
	GPU               string           `json:"gpu" form:"gpu"`             // GPU 数量/型号
	ClusterID         uint             `json:"clusterId" form:"clusterId"` // 集群ID
	EnableTensorboard bool             `json:"enableTensorboard" form:"enableTensorboard"`
	VolumeMounts      []VolumeMountReq `json:"volumeMounts" form:"volumeMounts"` // 挂载列表
	SSHKey            string           `json:"sshKey" form:"sshKey"`             // 公钥
	ExtraArgs         string           `json:"extraArgs" form:"extraArgs"`       // 启动参数
}

type VolumeMountReq struct {
	Name       string `json:"name"`       // 卷名称
	MountsPath string `json:"mountsPath"` // 容器内路径
	Size       int64  `json:"size"`       // 大小 (GB)
	Type       string `json:"type"`       // 类型 (如 dataset, model, workspace)
}

type UpdateNoteBookReq struct {
	Id     uint   `json:"id" form:"id"`
	Name   string `json:"name" form:"name" binding:"required"`
	Image  string `json:"image" form:"image"`
	CPU    string `json:"cpu" form:"cpu"`
	Memory string `json:"memory" form:"memory"`
}

type DeleteNoteBookReq struct {
	Id        uint   `json:"id" form:"id"`
	Name      string `json:"name" form:"name" binding:"required"`
	Namespace string `json:"namespace" form:"namespace"`
}

// ==================== Volume 文件存储 ====================

// 使用者信息
type VolumeUsage struct {
	Type string `json:"type"` // Notebook/TrainingJob
	Name string `json:"name"`
}

// 列表项
type VolumeItem struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`        // 存储名称
	Size        int64         `json:"size"`        // 容量
	Type        int           `json:"type"`        // 类型
	Status      string        `json:"status"`      // 未使用/使用中
	Area        string        `json:"area"`        // 区域
	CreatedAt   string        `json:"createdAt"`   // 创建时间
	UsedBy      []VolumeUsage `json:"usedBy"`      // 使用者列表
	ClusterId   uint          `json:"clusterId"`   // 集群ID
	ProductId   uint          `json:"productId"`   // 产品ID
	ProductName string        `json:"productName"` // 产品名称
}

// 列表响应
type VolumeListResp struct {
	List  []VolumeItem `json:"list"`
	Total int64        `json:"total"`
}

// 集群简要信息
type ClusterInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Area string `json:"area"`
}

// 区域/集群列表响应
type AreaListResp struct {
	Clusters []ClusterInfo `json:"clusters"`
}
