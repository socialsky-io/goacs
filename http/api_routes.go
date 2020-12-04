package http

import (
	"github.com/gin-gonic/gin"
	"goacs/http/controllers"
	"goacs/http/middleware/jwt"
	"goacs/lib"
)

func RegisterApiRoutes(gin *gin.Engine) {
	var env lib.Env
	gin.GET("/file/:filename", controllers.DownloadFile)
	apiGroup := gin.Group("/api")
	apiGroup.Use()
	apiGroup.POST("/auth/login", controllers.Login)

	apiGroup.Use(jwt.JWTAuthMiddleware(env.Get("JWT_SECRET", "")))
	{
		apiGroup.GET("/dashboard", controllers.GetDashboardData)
		apiGroup.POST("/user/create", controllers.UserCreate)

		apiGroup.GET("/device", controllers.GetDevicesList)
		apiGroup.GET("/device/:uuid", controllers.GetDevice)
		apiGroup.GET("/device/:uuid/kick", controllers.Kick)
		apiGroup.GET("/device/:uuid/parameters", controllers.GetDeviceParameters)
		apiGroup.POST("/device/:uuid/parameters", controllers.CreateParameter)
		apiGroup.PUT("/device/:uuid/parameters", controllers.UpdateParameter)
		apiGroup.POST("/device/:uuid/addobject", controllers.AddObject)
		apiGroup.POST("/device/:uuid/getparametervalues", controllers.GetParameterValues)
		apiGroup.GET("/device/:uuid/tasks", controllers.GetDeviceQueuedTasks)
		apiGroup.GET("/device/:uuid/templates", controllers.GetDeviceTemplates)

		apiGroup.GET("/template", controllers.GetTemplatesList)
		apiGroup.GET("/template/:templateid", controllers.GetTemplate)
		apiGroup.GET("/template/:templateid/parameters", controllers.GetTemplateParameters)

		apiGroup.GET("/faults/today", controllers.GetTodayFaults)

		apiGroup.POST("/file", controllers.UploadFile)
	}
}
