package initialize

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitViper() *viper.Viper {
	// 获取配置文件路径
	path, err := os.Executable()
	if err != nil {
		panic("初始化失败：可执行程序路径获取失败")
	}
	path = filepath.Dir(path)
	path = filepath.Join(path, "utils", "phoenix-config.yml")
	// 初始化viper
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	err = v.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		panic(err)
		//panic("初始化失败：读取配置文件失败")
	}
	return v
}
