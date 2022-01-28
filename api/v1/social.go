package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CalcExample godoc
// @Summary      calc example
// @Description  plus
// @Tags         example
// @Accept       json
// @Produce      json
// @Success      200   {integer}  string  "answer"
// @Failure      400   {string}   string  "ok"
// @Failure      404   {string}   string  "ok"
// @Failure      500   {string}   string  "ok"
// @Router       /api/v1/user [post]
func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "name": 1})
}
