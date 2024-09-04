package config

type Configuration struct {
	App          App    `mapstructure:"app" json:"app" yaml:"app"`
	Server       Server `mapstructure:"server" json:"server" yaml:"server"`
	Database     Database
	StaticServer StaticServer
	Template     Template
}
