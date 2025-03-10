package configs

type UploadConfig struct {
	MaxImageSize int                 `yaml:"maxImageSize" json:"maxImageSize"`
	MaxFileSize  int                 `yaml:"maxFileSize" json:"maxFileSize"`
	Uri          string              `yaml:"uri" json:"uri"`
	Path         string              `yaml:"path" json:"path"`
	Store        string              `yaml:"store" json:"store"`
	Github       *GithubUploadConfig `yaml:"github" json:"github"`
	VeymeToken   string              `yaml:"veymeToken" json:"veymeToken"`
}

type GithubUploadConfig struct {
	Token string `yaml:"token" json:"token"`
	User  string `yaml:"user" json:"user"`
	Repo  string `yaml:"repo" json:"repo"`
	Proxy string `yaml:"proxy" json:"proxy"`
}
