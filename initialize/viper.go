package initialize

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitViper() *viper.Viper {
	// 获取可执行文件所在目录
	rootPath, err := os.Executable()
	if err != nil {
		panic("初始化失败：可执行程序路径获取失败")
	}
	rootPath = filepath.Dir(rootPath)
	// 获取配置文件、题目、教程、用户头像保存路径
	path := filepath.Join(rootPath, "phoenix-config.yml")
	tutorialPath := filepath.Join(rootPath, "resource", "tutorials")
	problemPath := filepath.Join(rootPath, "resource", "problems")
	userPath := filepath.Join(rootPath, "resource", "users")
	codePath := filepath.Join(rootPath, "resource", "codes")
	// 创建资源文件夹
	err1 := os.MkdirAll(tutorialPath, os.ModePerm)
	err2 := os.MkdirAll(problemPath, os.ModePerm)
	err3 := os.MkdirAll(userPath, os.ModePerm)
	err4 := os.MkdirAll(codePath, os.ModePerm)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("初始化失败：初始化文件夹失败")
	}
	// 初始化viper，读取配置文件
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	err = v.ReadInConfig()
	if err != nil {
		panic("初始化失败：读取配置文件失败")
	}
	// 设置常用路径
	v.Set("root_path", rootPath)
	v.Set("problem_path", filepath.Join(rootPath, "resource", "problems"))
	v.Set("tutorial_path", filepath.Join(rootPath, "resource", "tutorials"))
	v.Set("user_path", filepath.Join(rootPath, "resource", "users"))
	v.Set("code_path", filepath.Join(rootPath, "resource", "codes"))
	return v
}
