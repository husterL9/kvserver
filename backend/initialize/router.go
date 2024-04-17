package initialize

import (
	"backend/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	// Router.Use(cors.New(cors.Config{
	// 	AllowAllOrigins: true, // 允许所有源
	// 	// AllowOrigins:     []string{"https://foo.com"},                               // 指定允许的源
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}, // 允许的HTTP方法
	// 	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},      // 允许的头部
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))
	systemRouter := router.RouterGroupApp.System

	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}

	{
		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
	}

	return Router
}
