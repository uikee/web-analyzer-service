package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/uikee/web-analyzer-service/internal/handler"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/ping", handler.PingHandler)
}
