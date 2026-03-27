package config

// SSHPiper SSHPiper 配置
type SSHPiper struct {
	// Host SSHPiper 服务的主机地址
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	// Port SSHPiper 服务的端口
	Port int `mapstructure:"port" json:"port" yaml:"port"`
}
