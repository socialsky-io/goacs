package http

import (
	"github.com/gin-gonic/gin"
	"goacs/http/controllers"
	"goacs/http/middleware/jwt"
	"goacs/lib"
)

func RegisterApiRoutes(gin *gin.Engine) {
	var env lib.Env
	apiGroup := gin.Group("/api")
	apiGroup.POST("/auth/login", controllers.Login)

	apiGroup.Use(jwt.JWTAuthMiddleware(env.Get("JWT_SECRET", "")))
	{
		apiGroup.GET("/dashboard", controllers.GetDashboardData)
		apiGroup.POST("/user/create", controllers.UserCreate)

		apiGroup.GET("/device", controllers.GetDevicesList)
		apiGroup.GET("/device/:uuid", controllers.GetDevice)
		apiGroup.GET("/device/:uuid/parameters", controllers.GetDeviceParameters)
		apiGroup.POST("/device/:uuid/parameters", controllers.CreateParameter)
		apiGroup.PUT("/device/:uuid/parameters", controllers.UpdateParameter)
	}
}
