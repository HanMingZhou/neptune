package piper

import (
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/pipe/request"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

type PiperApi struct{}

// AddPipe 创建Pipe
func (a *PiperApi) AddPipe(c *gin.Context) {
	var req request.AddPipeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}

	if err := piperService.CreatePipe(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "创建Pipe失败")
		return
	}
	response.OkWithMessage("创建Pipe成功", c)
}

// DeletePipe 删除Pipe
func (a *PiperApi) DeletePipe(c *gin.Context) {
	var req request.DeletePipeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if req.InstanceName == "" {
		response.FailWithMessage("Pipe名称不能为空", c)
		return
	}
	if err := piperService.DeletePipe(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "删除Pipe失败")
		return
	}
	response.OkWithMessage("删除Pipe成功", c)
}

// GetPipeList 获取Pipe列表
func (a *PiperApi) GetPipeList(c *gin.Context) {
	// TODO: 实现获取Pipe列表逻辑
	response.OkWithMessage("获取Pipe列表成功", c)
}

// UpdatePipe 更新Pipe
func (a *PiperApi) UpdatePipe(c *gin.Context) {
	// TODO: 实现更新Pipe逻辑
	response.OkWithMessage("更新Pipe成功", c)
}
