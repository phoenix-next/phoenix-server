package initialize

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitViper() *viper.Viper {
	// 获取配置文件路径
	path, err := os.Getwd()
	if err != nil {
		panic("初始化失败：工作目录获取失败")
	}
	path = filepath.Join(path, "phoenix-config.yml")
	// 初始化viper
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	err = v.ReadInConfig()
	if err != nil {
		panic("初始化失败：读取配置文件失败")
	}
	return v
}
