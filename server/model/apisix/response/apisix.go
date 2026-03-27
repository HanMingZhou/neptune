package response

type AuthApisixResp struct {
	UserID    uint   `json:"userId"`
	Namespace string `json:"namespace"`
	Token     string `json:"token,omitempty"`
	RateLimit int    `json:"rateLimit,omitempty"` // API-Key 限流配置
}
