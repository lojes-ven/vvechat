package router

import (
	"vvechat/internal/handler"
	"vvechat/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Launch() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/register", handler.Register)
		api.POST("/login/uid", handler.LoginByUid)
		api.POST("/login/phone_number", handler.LoginByPhone)
		auth := api.Group("/auth", middleware.JWTAuth())
		{
			auth.POST("/refresh_token", handler.RefreshToken)
		}
	}

	return r
}
