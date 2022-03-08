package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/global"
	"github.com/phoenix-next/phoenix-server/model"
	"runtime"
	"strconv"
)

func SolveUser(c *gin.Context) model.User {
	userRaw, _ := c.Get("user")
	return userRaw.(model.User)
}

func BindJsonData(c *gin.Context, model interface{}) interface{} {
	if err := c.ShouldBindJSON(&model); err != nil {
		_, file, line, _ := runtime.Caller(1)
		global.LOG.Panic(file + "(line " + strconv.Itoa(line) + "): bind model error")
	}
	return model
}
