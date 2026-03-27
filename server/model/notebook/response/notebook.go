package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func NoAuth(message string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		7,
		nil,
		message,
	})
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

// NotebookItem Notebook列表项
type NotebookItem struct {
	ID                uint   `json:"id"`                // Notebook ID
	DisplayName       string `json:"displayName"`       // Notebook名称
	InstanceName      string `json:"instanceName"`      // Notebook实例名称
	Namespace         string `json:"namespace"`         // Notebook所属命名空间
	Image             string `json:"image"`             // Notebook镜像
	CPU               int64  `json:"cpu"`               // Notebook CPU资源
	Memory            int64  `json:"memory"`            // Notebook 内存资源
	GPU               int64  `json:"gpu,omitempty"`     // Notebook GPU资源
	Status            string `json:"status"`            // Notebook 状态
	CreationTimestamp string `json:"creationTimestamp"` // Notebook 创建时间
	SSHCommand        string `json:"sshCommand"`        // SSH连接命令（密码登录）
	SSHKeyCommand     string `json:"sshKeyCommand"`     // SSH连接命令（公钥登录）
	SSHPassword       string `json:"sshPassword"`       // SSH登录密码（启用密码登录时返回）
	JupyterUrl        string `json:"jupyterUrl"`        // Jupyter访问地址
	TensorboardUrl    string `json:"tensorboardUrl"`    // TensorBoard访问地址（启用TensorBoard时返回）
	// 详情页补充字段
	ImageName          string      `json:"imageName"`          // 镜像展示名称
	GPUCount           int64       `json:"gpuCount"`           // GPU 数量
	GPUModel           string      `json:"gpuModel"`           // GPU 型号
	PayType            int         `json:"payType"`            // 计费方式
	Price              float64     `json:"price"`              // 单价
	EnableTensorboard  bool        `json:"enableTensorboard"`  // 是否启用 TensorBoard
	TensorboardLogPath string      `json:"tensorboardLogPath"` // TensorBoard 日志路径
	VolumeMounts       interface{} `json:"volumeMounts"`       // 挂载卷
	CreatedAt          string      `json:"createdAt"`          // 创建时间
}

// GetNotebookListResp 获取Notebook列表响应
type GetNotebookListResp struct {
	List  []NotebookItem `json:"list"`
	Total int64          `json:"total"`
}

// AddNotebookResp 创建Notebook响应
type AddNotebookResp struct {
	DisplayName string `json:"displayName"`
}

// UpdateNotebookResp 更新Notebook响应
type UpdateNotebookResp struct {
	DisplayName string `json:"displayName"`
}

// DeleteNotebookResp 删除Notebook响应
type DeleteNotebookResp struct {
	DisplayName string `json:"displayName"`
}

// PodInfoResp Pod 信息响应
type PodInfoResp struct {
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	Status     string   `json:"status"`
	HostIP     string   `json:"hostIP"`
	PodIP      string   `json:"podIP"`
	Containers []string `json:"containers"`
}

// NotebookAuthResp 认证结果（供 service 层返回给 api 层）
type NotebookAuthResp struct {
	UserID    uint   // 用户 ID
	Namespace string // 用户 Namespace
	Token     string // 原始 Token（用于设置 Cookie）
}

// TerminalInfoResp 终端连接所需信息（service 层返回给 api 层）
type TerminalInfoResp struct {
	Namespace    string
	PodName      string
	Container    string
	ClusterID    uint
	InstanceName string
}
