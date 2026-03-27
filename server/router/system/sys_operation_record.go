package system

import (
	"github.com/gin-gonic/gin"
)

type OperationRecordRouter struct{}

func (s *OperationRecordRouter) InitSysOperationRecordRouter(Router *gin.RouterGroup) {
	operationRecordRouter := Router.Group("operation/record")
	{

		operationRecordRouter.POST("delete", operationRecordApi.DeleteSysOperationRecord)            // 删除SysOperationRecord
		operationRecordRouter.POST("delete/multi", operationRecordApi.DeleteSysOperationRecordByIds) // 批量删除SysOperationRecord
		operationRecordRouter.GET("get", operationRecordApi.FindSysOperationRecord)                  // 根据ID获取SysOperationRecord
		operationRecordRouter.GET("list", operationRecordApi.GetSysOperationRecordList)              // 获取SysOperationRecord列表

	}
}
