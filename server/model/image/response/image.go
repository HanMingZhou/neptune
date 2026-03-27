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

// ImageItem 镜像列表项
type ImageItem struct {
	ID         uint   `json:"id"`         // 镜像 ID
	Name       string `json:"name"`       // 镜像名称
	Image      string `json:"image"`      // 镜像地址
	Size       string `json:"size"`       // 镜像大小
	ImageUUID  string `json:"imageUUID"`  // 镜像UUID
	Area       string `json:"area"`       // 镜像区域
	CreateTime string `json:"createTime"` // 镜像创建时间
	UsageType  int64  `json:"usageType"`  // 1 容器实例，2 训练镜像，3 推理镜像
	Type       int64  `json:"type"`       // 1 系统镜像，2 自定义镜像
	ImagePath  string `json:"imagePath"`  // 镜像路径
	UserId     uint   `json:"userId"`     // 用户ID
}

// GetImageListResp 获取镜像列表响应
type GetImageListResp struct {
	List  []ImageItem `json:"list"`
	Total int64       `json:"total"`
}

// AddImageResp 创建镜像响应
type AddImageResp struct {
	Name string `json:"name"`
}

// DeleteImageResp 删除镜像响应
type DeleteImageResp struct {
	Success bool `json:"success"`
}
