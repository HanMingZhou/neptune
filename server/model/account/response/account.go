package response

type AccessLogInfo struct {
	ID       uint   `json:"id"`
	Time     string `json:"time"`
	Ip       string `json:"ip"`
	Location string `json:"location"`
	Device   string `json:"device"`
	Os       string `json:"os"`
	Browser  string `json:"browser"`
	Method   string `json:"method"`
	Status   string `json:"status"`
	Reason   string `json:"reason"`
	Current  bool   `json:"current"` // 是否为当前会话
	LogType  string `json:"logType"`
}

type GetAccessLogListResp struct {
	List  []AccessLogInfo `json:"list"`
	Total int64           `json:"total"`
}

type SecurityItem struct {
	StatusText string `json:"statusText"`
	ActionText string `json:"actionText"`
	CanAction  bool   `json:"canAction"`
}

type SecurityStatusResp struct {
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	MfaEnabled     bool   `json:"mfaEnabled"`
	GithubBound    bool   `json:"githubBound"`
	GithubUsername string `json:"githubUsername"`
	AccessKeyId    string `json:"accessKeyId"`
	// UI渲染所需的派生状态信息
	PhoneStatus     string `json:"phoneStatus"`
	EmailStatus     string `json:"emailStatus"`
	MfaStatus       string `json:"mfaStatus"`
	GithubStatus    string `json:"githubStatus"`
	AccessKeyStatus string `json:"accessKeyStatus"`

	LastLoginTime string `json:"lastLoginTime"`
	LastLoginIp   string `json:"lastLoginIp"`
	SecurityScore int    `json:"securityScore"` // 计算得出的分数
}

type MfaSecretResp struct {
	Secret string `json:"secret"`
	Qr     string `json:"qr"` // Base64 图片或 URL
}

type GenerateAccessKeyResp struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}
