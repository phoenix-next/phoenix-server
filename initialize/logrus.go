package initialize

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Logger {
	// 获取日志路径
	path, err := os.Executable()
	if err != nil {
		panic("初始化失败：可执行程序路径获取失败")
	}
	path = filepath.Dir(path)
	path = filepath.Join(path, "phoenix-server.log")
	// 打开日志文件
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0666)
	if err != nil {
		panic("初始化失败：打开Server日志文件失败")
	}
	// 初始化logrus
	logger := logrus.New()
	logger.SetOutput(file)
	return logger
}
