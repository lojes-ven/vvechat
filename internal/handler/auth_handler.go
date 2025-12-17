package handler

import (
	"vvechat/internal/service"
	"vvechat/pkg/response"

	"github.com/gin-gonic/gin"
)

func RefreshToken(c *gin.Context) {
	id := c.GetUint64("id")
	typ := c.GetString("type")

	resp, err := service.RefreshToken(typ, id)
	if err != nil {
		response.Fail(c, 500, "token出现问题"+err.Error())
	}

	response.Success(c, "success", resp)
}
