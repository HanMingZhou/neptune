package tensorboard

import (
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/tensorboard/request"
	"gin-vue-admin/utils"

	"github.com/gin-gonic/gin"
)

type TensorBoardApi struct{}

// GetTensorBoardList 获取TensorBoard列表
func (a *TensorBoardApi) GetTensorBoardList(c *gin.Context) {

}

// AddTensorBoard 创建TensorBoard
func (a *TensorBoardApi) AddTensorBoard(c *gin.Context) {
	var req request.AddTensorBoardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}

	if err := tensorBoardService.CreateTensorboard(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "创建TensorBoard失败")
		return
	}
	response.OkWithMessage("创建TensorBoard成功", c)
}

// UpdateTensorBoard 更新TensorBoard
func (a *TensorBoardApi) UpdateTensorBoard(c *gin.Context) {

}

// DeleteTensorBoard 删除TensorBoard
func (a *TensorBoardApi) DeleteTensorBoard(c *gin.Context) {
	var req request.DeleteTensorBoardReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := tensorBoardService.DeleteTensorboard(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "删除TensorBoard失败")
		return
	}

	response.OkWithMessage("删除TensorBoard成功", c)
}
