package configs

type EmailConfig struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Host     string `yaml:"host" json:"host"`
	Addr     string `yaml:"addr" json:"addr"`
}
