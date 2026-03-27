package example

import (
	common "gin-vue-admin/model/common/request"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/example"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

type AttachmentCategoryApi struct{}

// GetCategoryList
// @Tags      GetCategoryList
// @Summary   媒体库分类列表
// @Security  AttachmentCategory
// @Produce   application/json
// @Success   200   {object}  response.Response{data=example.ExaAttachmentCategory,msg=string}  "媒体库分类列表"
// @Router    /attachmentCategory/getCategoryList [get]
func (a *AttachmentCategoryApi) GetCategoryList(c *gin.Context) {
	res, err := attachmentCategoryService.GetCategoryList()
	if err != nil {
		utils.HandleError(c, err, "获取分类列表失败")
		return
	}
	response.OkWithData(res, c)
}

// AddCategory
// @Tags      AddCategory
// @Summary   添加媒体库分类
// @Security  AttachmentCategory
// @accept    application/json
// @Produce   application/json
// @Param     data  body      example.ExaAttachmentCategory  true  "媒体库分类数据"// @Success   200   {object}  response.Response{msg=string}   "添加媒体库分类"
// @Router    /attachmentCategory/addCategory [post]
func (a *AttachmentCategoryApi) AddCategory(c *gin.Context) {
	var req example.ExaAttachmentCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}

	if err := attachmentCategoryService.AddCategory(&req); err != nil {
		utils.HandleError(c, err, "创建/更新失败")
		return
	}
	response.OkWithMessage("创建/更新成功", c)
}

// DeleteCategory
// @Tags      DeleteCategory
// @Summary   删除分类
// @Security  AttachmentCategory
// @accept    application/json
// @Produce   application/json
// @Param     data  body      common.GetById                true  "分类id"
// @Success   200   {object}  response.Response{msg=string}  "删除分类"
// @Router    /attachmentCategory/deleteCategory [post]
func (a *AttachmentCategoryApi) DeleteCategory(c *gin.Context) {
	var req common.GetById
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	if req.ID == 0 {
		response.FailWithMessage("参数错误", c)
		return
	}

	if err := attachmentCategoryService.DeleteCategory(&req.ID); err != nil {
		utils.HandleError(c, err, "删除失败")
		return
	}

	response.OkWithMessage("删除成功", c)
}
