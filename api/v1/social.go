package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-next/phoenix-server/model/api"
	"github.com/phoenix-next/phoenix-server/model/database"
	"github.com/phoenix-next/phoenix-server/service"
)

// @Summary      注册
// @Description  注册新用户
// @Tags         社交模块
// @Accept       json
// @Produce      json
// @Param        data  body      api.RegisterQ  true  "用户名, 邮箱, 密码, 验证码"
// @Success      200   {object}  api.RegisterA  "用户注册账号,返回注册成功信息"
// @Router       /api/v1/user [post]
func Register(c *gin.Context) {
	var data api.RegisterQ
	err := c.ShouldBindJSON(&data)
	if err != nil {
		panic(err)
	}
	err = service.CreateUser(&database.User{Name: data.Name, Password: data.Password, Email: data.Email})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, api.RegisterA{Message: "创建用户成功"})
}
