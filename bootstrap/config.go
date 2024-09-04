// config/config.go
package bootstrap

import (
	"fmt"
	"hserver/global"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port      string
		Host      string
		LogLevel  string
		LogFormat string
	}
	Database struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
	}
	StaticServer struct {
		Path string
	}
}

func InitializeConfig() *viper.Viper {
	config := "config.json"

	if configEnv := os.Getenv("VIPER_CONFIG"); configEnv != "" {
		config = configEnv
	}
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: $s \n", err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changed:", in.Name)
		// 配置文件修改后，重新加载配置
		if err := v.Unmarshal(&global.App.Config); err != nil {
			fmt.Println(err)
		}
	})
	// 将配置文件加载到全局变量中
	if err := v.Unmarshal(&global.App.Config); err != nil {
		fmt.Println(err)
	}

	return v
}
