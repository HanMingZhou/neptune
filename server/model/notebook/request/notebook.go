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

// 获取列表
type GetNotebookListReq struct {
	PageInfo
	UserId      uint   `json:"userId" form:"userId"`
	DisplayName string `json:"displayName" form:"displayName"`
	Type        string `json:"type" form:"type"`
	Status      string `json:"status" form:"status"`
	NotebookId  uint   `json:"notebookId" form:"notebookId"`
}

type AddNoteBookReq struct {
	UserId             uint                      `json:"userId" form:"userId"`                              // 用户ID
	DisplayName        string                    `json:"displayName" form:"displayName" binding:"required"` // Notebook名称
	InstanceName       string                    `json:"instanceName" form:"instanceName"`                  // Notebook实例名称 自动生成
	ImageId            uint                      `json:"imageId" form:"imageId"`                            // 镜像Id
	TensorBoard        bool                      `json:"tensorBoard" form:"tensorBoard"`                    // 是否启用Tensorboard
	TensorBoardLogPath string                    `json:"tensorBoardPath" form:"tensorBoardPath"`            // Tensorboard Log路径
	VolumeMounts       []*NotebookVolumeMountReq `json:"volumeMounts" form:"volumeMounts"`                  // 挂载列表
	SSHKeyId           uint                      `json:"sshKeyId" form:"sshKeyId"`                          // 公钥Id
	EnableSSHPassword  bool                      `json:"enableSshPassword" form:"enableSshPassword"`        // 是否启用SSH密码登录
	ExtraArgs          string                    `json:"extraArgs" form:"extraArgs"`                        // 启动参数
	PayType            int64                     `json:"payType" form:"payType"`                            // 计费方式  1:按量付费 2:包日 3:包周 4:包月
	Quantity           int64                     `json:"quantity" form:"quantity"`                          // 时长
	ProductId          uint                      `json:"productId" form:"productId"`                        // 产品ID
}

type NotebookVolumeMountReq struct {
	PVCId      uint   `json:"pvcId"`      // Volume ID
	MountsPath string `json:"mountsPath"` // 容器内路径
}

type UpdateNoteBookReq struct {
	Id          uint   `json:"id" form:"id"`
	UserId      uint   `json:"userId" form:"userId"`
	DisplayName string `json:"displayName" form:"displayName" binding:"required"`
	ImageId     uint   `json:"imageId" form:"imageId"`
	ChargeType  int64  `json:"chargeType" form:"chargeType"`
	Quantity    int64  `json:"quantity" form:"quantity"`
	ProductId   uint   `json:"productId" form:"productId"`
}

type DeleteNoteBookReq struct {
	Id     uint `json:"id" form:"id"`
	UserId uint `json:"userId" form:"userId"`
}

type GetNotebookDetailReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type GetNotebookLogsReq struct {
	ID         uint   `json:"id" form:"id" binding:"required"`
	Container  string `json:"container" form:"container"`
	Follow     bool   `json:"follow" form:"follow"`
	TailLines  int64  `json:"tailLines" form:"tailLines"`
	Timestamps bool   `json:"timestamps" form:"timestamps"`
}

// HandleTerminalReq Web Terminal 连接请求
type HandleTerminalReq struct {
	ID        uint   `form:"id" binding:"required"`
	Container string `form:"container"`
	Token     string `form:"token" binding:"required"`
}

// GetNotebookPodsReq 获取 Notebook Pod 列表请求
type GetNotebookPodsReq struct {
	ID uint `form:"id" binding:"required"`
}

// StopNotebookReq 停止 Notebook 请求
type StopNotebookReq struct {
	ID uint `json:"id" binding:"required"`
}

// StartNotebookReq 启动 Notebook 请求
type StartNotebookReq struct {
	ID uint `json:"id" binding:"required"`
}

// NotebookAuthReq Apisix forward-auth 认证请求（从 Header 中获取）
type NotebookAuthReq struct {
	OriginalUri string // X-Forwarded-Uri 或 X-Original-URI
	Token       string // 从 Cookie/Header/URL 参数中提取
}
