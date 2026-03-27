package request

// GetDashboardDataReq 获取仪表盘数据请求
type GetDashboardDataReq struct {
	UserId uint `json:"-"` // 从 JWT 获取，不由前端传入
	Days   int  `json:"days" form:"days"`
}
