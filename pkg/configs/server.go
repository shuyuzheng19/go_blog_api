package configs

// ServerConfig 服务器配置
type ServerConfig struct {
	Cron         bool       `yaml:"cron" json:"cron"`                 //是否开启定时任务
	Addr         string     `yaml:"addr" json:"addr"`                 //要监听的IP端口
	ApiPrefix    string     `yaml:"apiPrefix" json:"apiPrefix"`       //全局API的前缀
	ReadTimeOut  int        `yaml:"readTimeOut" json:"readTimeOut"`   //读取超时时间
	WriteTimeOut int        `yaml:"writeTimeOut" json:"writeTimeOut"` //写入超时时间
	Name         string     `yaml:"name" json:"name"`                 //APP名称
	MaxSize      int        `yaml:"maxSize" json:"maxSize"`           //请求体最大大小
	Cors         CorsConfig `yaml:"cors" json:"-"`
	Env          string     `yaml:"env"`
}

type CorsConfig struct {
	AllOrigins       bool     `yaml:"allOrigins"`
	Enable           bool     `yaml:"enable"`
	AllowOrigins     []string `yaml:"allowOrigins"`
	AllowMethods     []string `yaml:"allowMethods"`
	AllowHeaders     []string `yaml:"allowHeaders"`
	ExposeHeaders    []string `yaml:"exposeHeaders"`
	AllowCredentials bool     `yaml:"allowCredentials"`
}
