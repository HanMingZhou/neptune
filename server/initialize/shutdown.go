package initialize

import (
	"gin-vue-admin/global"
	"os"
	"os/signal"
	"syscall"

	"github.com/zeromicro/go-zero/core/logx"
)

// RegisterCloseSignal 注册关闭信号监听
// 在 main 函数中调用，用于优雅关闭系统资源
func RegisterCloseSignal() {
	// 1. 设置 defer 处理 panic 时的资源释放
	// 注意：这里的 defer 只能捕获调用 RegisterCloseSignal 的 goroutine 的 panic
	// 如果是在 main 中调用，它会在 main 退出前执行（前提是 main 正常返回，而不是 os.Exit）
	// 但由于我们使用了 signal.Notify 和 os.Exit，所以这里的 defer 主要是为了代码结构的完整性
	// 实际的关闭逻辑主要由信号处理 goroutine 触发

	// 2. 启动信号监听 goroutine
	go func() {
		sigCh := make(chan os.Signal, 1)
		// 监听 SIGINT (Ctrl+C) 和 SIGTERM (kill)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		// 阻塞等待信号
		sig := <-sigCh
		logx.Info("接收到退出信号，开始清理资源...", logx.Field("signal", sig.String()))

		// 执行清理逻辑
		CleanUp()

		logx.Info("资源清理完成，系统退出")
		os.Exit(0)
	}()
}

// CleanUp 执行所有资源清理工作
func CleanUp() {
	// 关闭 K8s Informer
	if global.GVA_K8S_CLIENT_INFO != nil && global.GVA_K8S_CLIENT_INFO.StopCh != nil {
		select {
		case <-global.GVA_K8S_CLIENT_INFO.StopCh:
			// 已经关闭，不做处理
		default:
			close(global.GVA_K8S_CLIENT_INFO.StopCh)
			logx.Info("K8s Informer stop channel closed")
		}
	}

	// 这里可以添加其他资源的清理逻辑，例如数据库连接关闭等
	// if global.GVA_DB != nil { ... }
}
