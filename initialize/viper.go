package initialize

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitViper() *viper.Viper {
	// 获取配置文件路径
	rootPath, err := os.Executable()
	if err != nil {
		panic("初始化失败：可执行程序路径获取失败")
	}
	rootPath = filepath.Dir(rootPath)
	path := filepath.Join(rootPath, "phoenix-config.yml")
	// 初始化viper
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	err = v.ReadInConfig()
	v.Set("root_path", rootPath)
	if err != nil {
		panic("初始化失败：读取配置文件失败")
	}
	return v
}
