package pvc

import (
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/pvc/request"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

type PVCApi struct{}

// AddPVC 创建PVC
func (a *PVCApi) AddPVC(c *gin.Context) {
	var req request.AddPVCReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}

	if err := pvcService.CreatePVCs(c.Request.Context(), &req, "pvc"); err != nil {
		utils.HandleError(c, err, "创建PVC失败")
		return
	}
	response.OkWithMessage("创建PVC成功", c)
}

// DeletePVC 删除PVC
func (a *PVCApi) DeletePVC(c *gin.Context) {
	var req request.DeletePVCReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if req.InstanceName == "" {
		response.FailWithMessage("PVC名称不能为空", c)
		return
	}
	if err := pvcService.DeletePVCs(c.Request.Context(), &req, "pvc"); err != nil {
		utils.HandleError(c, err, "删除PVC失败")
		return
	}
	response.OkWithMessage("删除PVC成功", c)
}

// GetPVCList 获取PVC列表
func (a *PVCApi) GetPVCList(c *gin.Context) {
	// TODO: 实现获取PVC列表逻辑
	response.OkWithMessage("获取PVC列表成功", c)
}

// UpdatePVC 更新PVC
func (a *PVCApi) UpdatePVC(c *gin.Context) {
	// TODO: 实现更新PVC逻辑
	response.OkWithMessage("更新PVC成功", c)
}
