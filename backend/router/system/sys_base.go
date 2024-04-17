package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("base")

	{
		baseRouter.POST("login", func(context *gin.Context) {
			context.JSON(http.StatusOK, "ok")
		})
		baseRouter.POST("register", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"message": "success",
			})
		})
	}

	return baseRouter
}
