package initialize

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/task"

	"github.com/robfig/cron/v3"
)

func Timer() {
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// 按量计费定时扣费（每小时执行）
		_, err = global.GVA_Timer.AddTaskByFunc("HourlyOrder", "@every 1h", func() {
			if err := task.HourlyOrderTask(); err != nil {
				fmt.Println("hourly order task error:", err)
			}
		}, "按量计费定时扣费", option...)
		if err != nil {
			fmt.Println("add hourly order timer error:", err)
		}

		// 存储按天计费（每天凌晨1点执行）
		_, err = global.GVA_Timer.AddTaskByFunc("DailyStorageOrder", "0 0 1 * * *", func() {
			if err := task.DailyStorageOrderTask(); err != nil {
				fmt.Println("daily storage order task error:", err)
			}
		}, "存储按天计费", option...)
		if err != nil {
			fmt.Println("add daily storage order timer error:", err)
		}
	}()
}
