package inference

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/consts"
	inferenceReq "gin-vue-admin/model/inference/request"
	"gin-vue-admin/service"
	terminalService "gin-vue-admin/service/terminal"
	"gin-vue-admin/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
)

type InferenceApi struct{}

var inferenceService = &service.ServiceGroupApp.InferenceServiceGroup.InferenceServiceService

// CreateInferenceService 创建推理服务
// @Summary 创建推理服务
// @Router /inference/services [post]
func (i *InferenceApi) CreateInferenceService(c *gin.Context) {
	var req inferenceReq.CreateInferenceServiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 从 JWT 获取用户 ID
	req.UserId = utils.GetUserID(c)

	resp, err := inferenceService.CreateInferenceService(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("创建推理服务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(resp, c)
}

// DeleteInferenceService 删除推理服务
// @Summary 删除推理服务
// @Router /inference/services [delete]
func (i *InferenceApi) DeleteInferenceService(c *gin.Context) {
	var req inferenceReq.DeleteInferenceServiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := inferenceService.DeleteInferenceService(c.Request.Context(), &req); err != nil {
		global.GVA_LOG.Error("删除推理服务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// StopInferenceService 停止推理服务
// @Summary 停止推理服务
// @Router /inference/services/stop [post]
func (i *InferenceApi) StopInferenceService(c *gin.Context) {
	var req inferenceReq.StopInferenceServiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := inferenceService.StopInferenceService(c.Request.Context(), &req); err != nil {
		global.GVA_LOG.Error("停止推理服务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("停止成功", c)
}

// StartInferenceService 启动推理服务
// @Summary 启动推理服务
// @Router /inference/services/start [post]
func (i *InferenceApi) StartInferenceService(c *gin.Context) {
	var req inferenceReq.StartInferenceServiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := inferenceService.StartInferenceService(c.Request.Context(), &req); err != nil {
		global.GVA_LOG.Error("启动推理服务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("启动成功", c)
}

// GetInferenceServiceList 获取推理服务列表
// @Summary 获取推理服务列表
// @Router /inference/services/list [post]
func (i *InferenceApi) GetInferenceServiceList(c *gin.Context) {
	var req inferenceReq.GetInferenceServiceListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	req.UserId = utils.GetUserID(c)

	resp, err := inferenceService.GetInferenceServiceList(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("获取推理服务列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(resp, c)
}

// GetInferenceServiceDetail 获取推理服务详情
// @Summary 获取推理服务详情
// @Router /inference/services/detail [get]
func (i *InferenceApi) GetInferenceServiceDetail(c *gin.Context) {
	var req inferenceReq.GetInferenceServiceDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	resp, err := inferenceService.GetInferenceServiceDetail(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("获取推理服务详情失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(resp, c)
}

// GetInferenceServiceLogs 获取推理服务日志
// @Summary 获取推理服务日志
// @Router /inference/services/logs [get]
func (i *InferenceApi) GetInferenceServiceLogs(c *gin.Context) {
	var req inferenceReq.GetInferenceServiceLogsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	reader, err := inferenceService.GetInferenceServiceLogs(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("获取推理服务日志失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer reader.Close()

	// 读取日志内容，以 JSON 格式返回（与前端 axios 拦截器兼容）
	data, readErr := io.ReadAll(reader)
	if readErr != nil {
		response.FailWithMessage("读取日志失败", c)
		return
	}
	response.OkWithData(string(data), c)
}

// GetInferenceServicePods 获取推理服务Pod列表
// @Summary 获取推理服务Pod列表
// @Router /inference/services/pods [get]
func (i *InferenceApi) GetInferenceServicePods(c *gin.Context) {
	var req inferenceReq.GetInferenceServicePodsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	pods, err := inferenceService.GetInferenceServicePods(c.Request.Context(), req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取Pod列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(pods, c)
}

// HandleStreamLogs 处理推理服务日志流 WebSocket 连接
// @Summary 推理服务日志流
// @Router /inference/log/stream [get]
func (i *InferenceApi) HandleStreamLogs(c *gin.Context) {
	// 1. 鉴权
	token := c.Query("token")
	if token == "" {
		response.NoAuth("未登录或非法访问", c)
		return
	}
	j := utils.NewJWT()
	if _, err := j.ParseToken(token); err != nil {
		response.NoAuth("授权已过期", c)
		return
	}

	// 2. 解析参数
	var req inferenceReq.GetInferenceServiceLogsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 3. 升级到 WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 4096,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logx.Error("WebSocket升级失败")
		return
	}
	defer conn.Close()

	// 4. 获取日志流
	req.Follow = true
	logStream, err := inferenceService.GetInferenceServiceLogs(c.Request.Context(), &req)
	if err != nil {
		logx.Error("获取日志流失败")
		conn.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
		return
	}
	defer logStream.Close()

	// 5. 监听断开
	done := make(chan struct{})
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				close(done)
				return
			}
		}
	}()

	// 6. 发送日志
	buf := make([]byte, 4096)
	for {
		select {
		case <-done:
			return
		default:
			n, err := logStream.Read(buf)
			if err != nil {
				if err != io.EOF {
					logx.Error("读取日志流失败")
				}
				conn.WriteMessage(websocket.TextMessage, []byte("\n--- 日志流结束 ---\n"))
				return
			}
			if n > 0 {
				if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					return
				}
			}
		}
	}
}

// HandleTerminal 处理推理服务 Web Terminal WebSocket 连接
// @Summary 推理服务终端连接
// @Router /inference/terminal [get]
func (i *InferenceApi) HandleTerminal(c *gin.Context) {
	// 1. 解析参数
	var req inferenceReq.HandleTerminalReq
	if err := c.ShouldBindQuery(&req); err != nil {
		logx.Error("参数错误", err)
		return
	}

	// 2. JWT 鉴权（WebSocket 无法携带 Header，通过 query 传 token）
	j := utils.NewJWT()
	if _, err := j.ParseToken(req.Token); err != nil {
		logx.Error("Token 无效")
		return
	}

	// 3. 获取服务信息
	svcInfo, err := inferenceService.GetInferenceServiceForTerminal(c.Request.Context(), req.ID)
	if err != nil {
		logx.Error("获取推理服务信息失败", err)
		return
	}

	// 4. 确定 Pod 名称
	podName := req.PodName
	if podName == "" {
		// 默认连接 head Pod
		if svcInfo.DeployType == consts.DeployTypeDistributed {
			podName = "" // 通过 label selector 查找
		}
	}

	container := req.Container
	if container == "" {
		if svcInfo.DeployType == consts.DeployTypeDistributed {
			container = svcInfo.InstanceName + "-head"
		} else {
			container = svcInfo.InstanceName
		}
	}

	labelSelector := fmt.Sprintf("%s=%s,neptune.io/instance=%s", consts.LabelApp, svcInfo.InstanceName, svcInfo.InstanceName)

	// 5. 升级到 WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logx.Error("WebSocket 升级失败")
		return
	}
	defer conn.Close()

	// 6. 使用通用 Terminal 服务
	termReq := &terminalService.TerminalRequest{
		Namespace:     svcInfo.Namespace,
		PodName:       podName,
		LabelSelector: labelSelector,
		Container:     container,
		ClusterID:     svcInfo.ClusterID,
		Command:       []string{"/bin/bash"},
	}
	if err := terminalService.TerminalServiceApp.HandleWebSocket(c.Request.Context(), termReq, conn); err != nil {
		logx.Error("推理服务 Terminal 会话结束")
		conn.WriteMessage(websocket.TextMessage, []byte("\r\n\033[31mError: "+err.Error()+"\033[0m\r\n"))
	}
}
