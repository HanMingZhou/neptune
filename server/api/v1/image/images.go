package image

import (
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/image/request"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

type ImageApi struct{}

// 获取镜像列表
func (a *ImageApi) GetImageList(c *gin.Context) {
	var req request.GetImageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}
	userId := utils.GetUserID(c)
	res, err := imageService.GetImageList(c.Request.Context(), &req, userId)
	if err != nil {
		utils.HandleError(c, err, "获取镜像列表失败")
		return
	}
	response.OkWithData(res, c)
}

// 删除镜像
func (a *ImageApi) DeleteImage(c *gin.Context) {
	var req request.DeleteImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	userId := utils.GetUserID(c)
	if err := imageService.DeleteImage(c.Request.Context(), &req, userId); err != nil {
		utils.HandleError(c, err, "删除失败")
		return
	}

	response.OkWithMessage("删除成功", c)
}

// 创建镜像
func (a *ImageApi) CreateImage(c *gin.Context) {
	var req request.AddImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	userId := utils.GetUserID(c)
	if err := imageService.CreateImage(c.Request.Context(), &req, userId); err != nil {
		utils.HandleError(c, err, "创建镜像失败")
		return
	}
	response.OkWithMessage("创建成功", c)
}

// 更新镜像
func (a *ImageApi) UpdateImage(c *gin.Context) {
	var req request.UpdateImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	userId := utils.GetUserID(c)
	if err := imageService.UpdateImage(c.Request.Context(), &req, userId); err != nil {
		utils.HandleError(c, err, "更新镜像失败")
		return
	}
	response.OkWithMessage("更新成功", c)
}
