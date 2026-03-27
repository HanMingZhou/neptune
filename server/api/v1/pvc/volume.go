package pvc

import (
	"gin-vue-admin/model/common/response"
	pvcReq "gin-vue-admin/model/pvc/request"
	"gin-vue-admin/service"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type VolumeApi struct{}

var volumeService = service.ServiceGroupApp.PVCServiceGroup.VolumeService

// CreateVolume 创建文件存储
// @Tags Volume
// @Summary 创建文件存储
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body pvcReq.CreateVolumeReq true "创建参数"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /api/v1/volume/create [post]
func (v *VolumeApi) CreateVolume(c *gin.Context) {
	var req pvcReq.CreateVolumeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}

	// 获取当前用户信息
	userId := utils.GetUserID(c)
	if userId == 0 {
		utils.HandleError(c, errors.New("请先登录"), "请先登录")
		return
	}
	namespace := utils.GetUserNamespace(c)
	req.Namespace = namespace
	req.UserId = userId
	if err := volumeService.CreateVolume(c.Request.Context(), &req); err != nil {
		logx.Error("创建存储失败")
		utils.HandleError(c, err, "")
		return
	}

	response.OkWithMessage("创建成功", c)
}

// GetVolumeList 获取文件存储列表
// @Tags Volume
// @Summary 获取文件存储列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query pvcReq.GetVolumeListReq true "查询参数"
// @Success 200 {object} response.Response{data=pvcResp.VolumeListResp} "获取成功"
// @Router /api/v1/volume/list [get]
func (v *VolumeApi) GetVolumeList(c *gin.Context) {
	var req pvcReq.GetVolumeListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleError(c, err, "参数错误:")
		return
	}

	userId := utils.GetUserID(c)
	if userId == 0 {
		utils.HandleError(c, errors.New("请先登录"), "请先登录")
		return
	}
	req.UserId = userId

	result, err := volumeService.GetVolumeList(c.Request.Context(), &req)
	if err != nil {
		logx.Error("获取存储列表失败")
		utils.HandleError(c, err, "")
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}

// ExpandVolume 扩容文件存储
// @Tags Volume
// @Summary 扩容文件存储
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body pvcReq.ExpandVolumeReq true "扩容参数"
// @Success 200 {object} response.Response{msg=string} "扩容成功"
// @Router /api/v1/volume/expand [post]
func (v *VolumeApi) ExpandVolume(c *gin.Context) {
	var req pvcReq.ExpandVolumeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误:")
		return
	}

	userId := utils.GetUserID(c)
	if userId == 0 {
		utils.HandleError(c, errors.New("请先登录"), "请先登录")
		return
	}
	req.UserId = userId

	if err := volumeService.ExpandVolume(c.Request.Context(), &req); err != nil {
		logx.Error("扩容存储失败")
		utils.HandleError(c, err, "")
		return
	}

	response.OkWithMessage("扩容成功", c)
}

// DeleteVolume 删除文件存储
// @Tags Volume
// @Summary 删除文件存储
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body pvcReq.DeleteVolumeReq true "删除参数"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /api/v1/volume/delete [post]
func (v *VolumeApi) DeleteVolume(c *gin.Context) {
	var req pvcReq.DeleteVolumeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误:")
		return
	}

	userId := utils.GetUserID(c)
	if userId == 0 {
		utils.HandleError(c, errors.New("请先登录"), "请先登录")
		return
	}
	req.UserId = userId

	if err := volumeService.DeleteVolume(c.Request.Context(), &req); err != nil {
		logx.Error("删除存储失败")
		utils.HandleError(c, err, "")
		return
	}

	response.OkWithMessage("删除成功", c)
}

// GetAreaList 获取可用区域列表
// @Tags Volume
// @Summary 获取可用区域列表
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=pvcResp.AreaListResp} "获取成功"
// @Router /api/v1/volume/areas [get]
func (v *VolumeApi) GetAreaList(c *gin.Context) {
	result, err := volumeService.GetAreaList(c.Request.Context())
	if err != nil {
		logx.Error("获取区域列表失败")
		utils.HandleError(c, err, "")
		return
	}

	response.OkWithDetailed(result, "获取成功", c)
}
