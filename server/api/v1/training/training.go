package training

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/consts"
	trainingModel "gin-vue-admin/model/training"
	"gin-vue-admin/model/training/request"
	terminalService "gin-vue-admin/service/terminal"
	trainingService "gin-vue-admin/service/training"
	"gin-vue-admin/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type TrainingJobApi struct{}

// CreateTrainingJob 创建训练任务
// @Tags TrainingJob
// @Summary 创建训练任务
// @Accept json
// @Produce json
// @Param data body request.CreateTrainingJobReq true "创建训练任务"
// @Success 200 {object} response.Response{}
// @Router /training/create [post]
func (api *TrainingJobApi) CreateTrainingJob(c *gin.Context) {
	var req request.CreateTrainingJobReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 获取用户信息
	userId := utils.GetUserID(c)
	namespace := utils.GetUserNamespace(c)
	req.UserId = userId
	req.Namespace = namespace

	resp, err := trainingService.TrainingJobServiceApp.CreateTrainingJob(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err, "创建训练任务失败")
		return
	}

	response.OkWithDetailed(resp, "创建成功", c)
}

// DeleteTrainingJob 删除训练任务
// @Tags TrainingJob
// @Summary 删除训练任务
// @Accept json
// @Produce json
// @Param data body request.DeleteTrainingJobReq true "删除训练任务"
// @Success 200 {object} response.Response{}
// @Router /training/delete [post]
func (api *TrainingJobApi) DeleteTrainingJob(c *gin.Context) {
	var req request.DeleteTrainingJobReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := trainingService.TrainingJobServiceApp.DeleteTrainingJob(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "删除训练任务失败")
		return
	}

	response.OkWithMessage("删除成功", c)
}

// StopTrainingJob 停止训练任务
// @Tags TrainingJob
// @Summary 停止训练任务
// @Accept json
// @Produce json
// @Param data body request.StopTrainingJobReq true "停止训练任务"
// @Success 200 {object} response.Response{}
// @Router /training/stop [post]
func (api *TrainingJobApi) StopTrainingJob(c *gin.Context) {
	var req request.StopTrainingJobReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := trainingService.TrainingJobServiceApp.StopTrainingJob(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "停止训练任务失败")
		return
	}

	response.OkWithMessage("停止成功", c)
}

// GetTrainingJobList 获取训练任务列表
// @Tags TrainingJob
// @Summary 获取训练任务列表
// @Accept json
// @Produce json
// @Param data query request.GetTrainingJobListReq true "获取训练任务列表"
// @Success 200 {object} response.Response{}
// @Router /training/list [get]
func (api *TrainingJobApi) GetTrainingJobList(c *gin.Context) {
	var req request.GetTrainingJobListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	req.UserId = userId

	resp, err := trainingService.TrainingJobServiceApp.GetTrainingJobList(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err, "获取训练任务列表失败")
		return
	}

	response.OkWithDetailed(resp, "获取成功", c)
}

// GetTrainingJobDetail 获取训练任务详情
// @Tags TrainingJob
// @Summary 获取训练任务详情
// @Accept json
// @Produce json
// @Param data query request.GetTrainingJobDetailReq true "获取训练任务详情"
// @Success 200 {object} response.Response{}
// @Router /training/get [get]
func (api *TrainingJobApi) GetTrainingJobDetail(c *gin.Context) {
	var req request.GetTrainingJobDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	resp, err := trainingService.TrainingJobServiceApp.GetTrainingJobDetail(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err, "获取训练任务详情失败")
		return
	}

	response.OkWithDetailed(resp, "获取成功", c)
}

// GetTrainingJobPods 获取训练任务 Pod 列表
// @Tags TrainingJob
// @Summary 获取训练任务 Pod 列表
// @Accept json
// @Produce json
// @Param id query int true "任务ID"
// @Success 200 {object} response.Response{}
// @Router /training/pod/list [get]
func (api *TrainingJobApi) GetTrainingJobPods(c *gin.Context) {
	var req request.GetTrainingJobPodsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	pods, err := trainingService.TrainingJobServiceApp.GetTrainingJobPods(c.Request.Context(), req.ID)
	if err != nil {
		utils.HandleError(c, err, "获取 Pod 列表失败")
		return
	}

	response.OkWithDetailed(pods, "获取成功", c)
}

// GetTrainingJobLogs 获取训练任务日志
// @Tags TrainingJob
// @Summary 获取训练任务日志
// @Accept json
// @Produce json
// @Param data query request.GetTrainingJobLogsReq true "获取训练任务日志"
// @Success 200 {object} response.Response{}
// @Router /training/log/list [get]
func (api *TrainingJobApi) GetTrainingJobLogs(c *gin.Context) {
	var req request.GetTrainingJobLogsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	logStream, err := trainingService.TrainingJobServiceApp.GetTrainingJobLogs(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err, "获取训练任务日志失败")
		return
	}
	defer logStream.Close()

	// 设置响应头
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Cache-Control", "no-cache")

	// 流式输出日志
	if req.Follow {
		c.Stream(func(w io.Writer) bool {
			buf := make([]byte, 4096)
			n, err := logStream.Read(buf)
			if err != nil {
				return false
			}
			w.Write(buf[:n])
			return true
		})
	} else {
		// 非流式，读取全部并返回 JSON
		content, err := io.ReadAll(logStream)
		if err != nil {
			utils.HandleError(c, err, "读取日志失败")
			return
		}
		response.OkWithDetailed(gin.H{"logs": string(content)}, "获取成功", c)
	}
}

// HandleTerminal 处理 Web Terminal WebSocket 连接
// @Tags TrainingJob
// @Summary Web Terminal (WebSocket)
// @Accept json
// @Produce json
// @Param id query int true "任务ID"
// @Param taskName query string false "Task名称"
// @Param container query string false "容器名称"
// @Router /training/terminal [get]
func (api *TrainingJobApi) HandleTerminal(c *gin.Context) {
	// 1. 手动鉴权 (因为 WebSocket 握手无法发送自定义 Header)
	token := c.Query("token")
	if token == "" {
		response.NoAuth("未登录或非法访问", c)
		return
	}
	j := utils.NewJWT()
	_, err := j.ParseToken(token)
	if err != nil {
		response.NoAuth("授权已过期", c)
		return
	}

	// 2. 解析参数
	var req request.TerminalReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 2. 获取任务信息
	var job trainingModel.TrainingJob
	if err := global.GVA_DB.Where("id = ?", req.ID).First(&job).Error; err != nil {
		utils.HandleError(c, err, "任务不存在")
		return
	}

	// 3. 检查任务状态
	if job.Status != consts.TrainingStatusRunning && job.Status != consts.TrainingStatusPending {
		response.FailWithMessage("任务未在运行中，无法连接终端", c)
		return
	}

	// 4. 构建标签选择器
	// Volcano Job 的 Pod 标签通常包含 volcano.sh/job-name
	labelSelector := fmt.Sprintf("%s=%s", consts.LabelVolcanoJobName, job.JobName)
	if req.TaskName != "" {
		labelSelector += fmt.Sprintf(",%s=%s", consts.LabelVolcanoTaskSpec, req.TaskName)
	}

	// 5. 确定容器名称
	containerName := req.Container
	if containerName == "" {
		containerName = job.JobName + "-" + req.TaskName
	}

	// 6. 升级到 WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logx.Error("WebSocket升级失败")
		return
	}
	defer func() {
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Terminal session ended"))
		conn.Close()
	}()

	// 7. 调用通用 Terminal 服务
	termReq := &terminalService.TerminalRequest{
		Namespace:     job.Namespace,
		LabelSelector: labelSelector,
		Container:     containerName,
		ClusterID:     job.ClusterID,
	}

	if err := terminalService.TerminalServiceApp.HandleWebSocket(c.Request.Context(), termReq, conn); err != nil {
		logx.Error("Terminal会话结束")
		conn.WriteMessage(websocket.TextMessage, []byte("\r\n\033[31mError: "+err.Error()+"\033[0m\r\n"))
	}
}

// HandleStreamLogs 处理日志流 WebSocket 连接
// @Tags TrainingJob
// @Summary 日志流 (WebSocket)
// @Accept json
// @Produce json
// @Param id query int true "任务ID"
// @Param taskName query string false "Task名称"
// @Param podIndex query int false "Pod索引"
// @Router /training/log/stream [get]
func (api *TrainingJobApi) HandleStreamLogs(c *gin.Context) {
	// 1. 手动鉴权 (因为 WebSocket 握手无法发送自定义 Header)
	token := c.Query("token")
	if token == "" {
		response.NoAuth("未登录或非法访问", c)
		return
	}
	j := utils.NewJWT()
	_, err := j.ParseToken(token)
	if err != nil {
		response.NoAuth("授权已过期", c)
		return
	}

	// 2. 解析参数
	var req request.GetTrainingJobLogsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
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

	// 4. 设置 Follow 为 true，获取流式日志
	req.Follow = true
	logStream, err := trainingService.TrainingJobServiceApp.GetTrainingJobLogs(c.Request.Context(), &req)
	if err != nil {
		logx.Error("获取日志流失败")
		conn.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
		return
	}
	defer logStream.Close()

	// 5. 启动一个 goroutine 监听客户端断开
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

	// 6. 流式读取日志并发送到 WebSocket
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
				// 发送结束标记
				conn.WriteMessage(websocket.TextMessage, []byte("\n--- 日志流结束 ---\n"))
				return
			}
			if n > 0 {
				if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					logx.Error("发送日志失败")
					return
				}
			}
		}
	}
}
