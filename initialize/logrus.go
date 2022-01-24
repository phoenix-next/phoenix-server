package initialize

import "github.com/sirupsen/logrus"

func InitLogger() *logrus.Logger {
	logger := logrus.New()
	return logger
}
