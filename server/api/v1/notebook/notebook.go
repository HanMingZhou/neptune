package notebook

import (
	"gin-vue-admin/model/common/response"
	"gin-vue-admin/model/notebook/request"
	terminalService "gin-vue-admin/service/terminal"
	"gin-vue-admin/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type NoteBookApi struct{}

// GetNotebookList 获取列表
func (a *NoteBookApi) GetNotebookList(c *gin.Context) {
	var req request.GetNotebookListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}
	userId := utils.GetUserID(c)
	req.UserId = userId
	res, err := noteBookService.GetNotebookList(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err, "获取分类列表失败")
		return
	}
	response.OkWithData(res, c)
}

// AddNoteBook 创建 Notebook
func (a *NoteBookApi) AddNoteBook(c *gin.Context) {
	var req request.AddNoteBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}
	userId := utils.GetUserID(c)
	req.UserId = userId
	if err := noteBookService.CreateNotebook(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "创建/更新失败")
		return
	}
	response.OkWithMessage("创建/更新成功", c)
}

// UpdateNoteBook 更新 Notebook
func (a *NoteBookApi) UpdateNoteBook(c *gin.Context) {
	var req request.UpdateNoteBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}
	userId := utils.GetUserID(c)
	req.UserId = userId
	if err := noteBookService.UpdateNotebook(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "更新失败")
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteNoteBook 删除 Notebook
func (a *NoteBookApi) DeleteNoteBook(c *gin.Context) {
	var req request.DeleteNoteBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}
	userId := utils.GetUserID(c)
	req.UserId = userId
	if err := noteBookService.DeleteNotebook(c.Request.Context(), &req); err != nil {
		utils.HandleError(c, err, "删除失败")
		return
	}
	response.OkWithMessage("删除成功", c)
}

// HandleTerminal 处理 Web Terminal WebSocket 连接
func (a *NoteBookApi) HandleTerminal(c *gin.Context) {
	// 1. 解析参数（token 鉴权通过 req 绑定）
	var req request.HandleTerminalReq
	if err := c.ShouldBindQuery(&req); err != nil {
		logx.Error("参数错误", err)
		return
	}

	// 2. 鉴权
	j := utils.NewJWT()
	if _, err := j.ParseToken(req.Token); err != nil {
		logx.Error("Token 无效")
		return
	}

	// 3. 从 service 层获取终端连接信息
	termInfo, err := noteBookService.GetTerminalInfo(c.Request.Context(), &req)
	if err != nil {
		logx.Error("获取终端信息失败", err)
		return
	}

	// 4. 升级到 WebSocket
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

	// 5. 使用通用 Terminal 服务
	termReq := &terminalService.TerminalRequest{
		Namespace: termInfo.Namespace,
		PodName:   termInfo.PodName,
		Container: termInfo.Container,
		ClusterID: termInfo.ClusterID,
		Command:   []string{"/bin/bash"},
	}
	if err := terminalService.TerminalServiceApp.HandleWebSocket(c.Request.Context(), termReq, conn); err != nil {
		logx.Error("Notebook Terminal会话结束")
		conn.WriteMessage(websocket.TextMessage, []byte("\r\n\033[31mError: "+err.Error()+"\033[0m\r\n"))
	}
}

// GetNotebookDetail 获取 Notebook 详情
func (a *NoteBookApi) GetNotebookDetail(c *gin.Context) {
	var req request.GetNotebookDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}
	resp, err := noteBookService.GetNotebookDetail(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err, "获取详情失败")
		return
	}
	response.OkWithDetailed(resp, "获取成功", c)
}

// GetNotebookPods 获取 Notebook Pod 列表
func (a *NoteBookApi) GetNotebookPods(c *gin.Context) {
	var req request.GetNotebookPodsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}
	pods, err := noteBookService.GetNotebookPods(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err, "获取 Pod 列表失败")
		return
	}
	response.OkWithDetailed(pods, "获取成功", c)
}

// GetNotebookLogs 获取 Notebook 日志
func (a *NoteBookApi) GetNotebookLogs(c *gin.Context) {
	var req request.GetNotebookLogsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}
	logStream, err := noteBookService.GetNotebookLogs(c.Request.Context(), &req)
	if err != nil {
		utils.HandleError(c, err, "获取日志失败")
		return
	}
	defer logStream.Close()

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Header("Cache-Control", "no-cache")

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
		content, err := io.ReadAll(logStream)
		if err != nil {
			utils.HandleError(c, err, "读取日志失败")
			return
		}
		response.OkWithDetailed(gin.H{"logs": string(content)}, "获取成功", c)
	}
}

// HandleStreamLogs 处理日志流 WebSocket 连接
func (a *NoteBookApi) HandleStreamLogs(c *gin.Context) {
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
	var req request.GetNotebookLogsReq
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

	// 4. 获取日志流
	req.Follow = true
	logStream, err := noteBookService.GetNotebookLogs(c.Request.Context(), &req)
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

// StopNotebook 停止 Notebook
func (a *NoteBookApi) StopNotebook(c *gin.Context) {
	var req request.StopNotebookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}
	if err := noteBookService.StopNotebook(c.Request.Context(), req.ID); err != nil {
		utils.HandleError(c, err, "停止失败")
		return
	}
	response.OkWithMessage("停止成功", c)
}

// StartNotebook 启动 Notebook
func (a *NoteBookApi) StartNotebook(c *gin.Context) {
	var req request.StartNotebookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, err, "参数错误")
		return
	}
	if err := noteBookService.StartNotebook(c.Request.Context(), req.ID); err != nil {
		utils.HandleError(c, err, "启动失败")
		return
	}
	response.OkWithMessage("启动成功", c)
}
