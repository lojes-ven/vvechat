package router

import (
	"vvechat/internal/handler"
	"vvechat/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Launch(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	serve := service.NewUserService(db)
	handle := handler.NewUserHandler(serve)

	api := r.Group("/api")
	{
		api.POST("/register", handle.Register)
		api.POST("/login/uid", handle.LoginByUid)
	}

	return r
}
