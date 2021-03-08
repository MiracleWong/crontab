package common

// 定时任务
type Job struct {
	Name int `json:"name"`			// 任务名
	Command int `json:"command"`	// shell 命令
	CronExpr int `json:"cronExpr"`	// cron 表达式
}


