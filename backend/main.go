package main

import (
	"relay-panel/backend/controllers"
	"relay-panel/backend/database"
	"relay-panel/backend/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	public := r.Group("/api")
	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)
	public.POST("/flow/upload", controllers.UploadFlowData)
	public.GET("/captcha/generate", controllers.GenerateCaptcha)
	public.GET("/captcha/:captchaId", controllers.ServeCaptcha)
	public.POST("/captcha/verify", controllers.VerifyCaptcha)

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/users", controllers.GetUsers)
		protected.PUT("/users/:id", controllers.UpdateUser)
		protected.DELETE("/users/:id", controllers.DeleteUser)
		protected.POST("/users/update_password", controllers.UpdatePassword)
		protected.POST("/users/:id/reset_traffic", controllers.ResetTraffic)
		protected.POST("/nodes", controllers.CreateNode)
		protected.GET("/nodes", controllers.GetNodes)
		protected.GET("/nodes/:id", controllers.GetNode)
		protected.POST("/nodes/:id/install", controllers.GetInstallCommand)
		protected.PUT("/nodes/:id", controllers.UpdateNode)
		protected.DELETE("/nodes/:id", controllers.DeleteNode)

		protected.POST("/tunnels", controllers.CreateTunnel)
		protected.GET("/tunnels", controllers.GetTunnels)
		protected.GET("/tunnels/:id", controllers.GetTunnel)
		protected.PUT("/tunnels/:id", controllers.UpdateTunnel)
		protected.DELETE("/tunnels/:id", controllers.DeleteTunnel)
		protected.POST("/tunnels/assign", controllers.AssignUserTunnel)
		protected.GET("/users/:id/tunnels", controllers.GetUserTunnels)
		protected.DELETE("/users/:id/tunnels/:tunnel_id", controllers.RemoveUserTunnel)
		protected.GET("/tunnels/:id/diagnose", controllers.DiagnoseTunnel)

		protected.POST("/forwards", controllers.CreateForward)
		protected.GET("/forwards", controllers.GetForwards)
		protected.GET("/forwards/:id", controllers.GetForward)
		protected.PUT("/forwards/:id", controllers.UpdateForward)
		protected.DELETE("/forwards/:id", controllers.DeleteForward)
		protected.POST("/forwards/:id/pause", controllers.PauseForward)
		protected.POST("/forwards/:id/resume", controllers.ResumeForward)
		protected.GET("/forwards/:id/diagnose", controllers.DiagnoseForward)
		protected.POST("/forwards/reorder", controllers.ReorderForward)

		protected.POST("/speedlimits", controllers.CreateSpeedLimit)
		protected.GET("/speedlimits", controllers.GetSpeedLimits)
		protected.GET("/speedlimits/:id", controllers.GetSpeedLimit)
		protected.PUT("/speedlimits/:id", controllers.UpdateSpeedLimit)
		protected.DELETE("/speedlimits/:id", controllers.DeleteSpeedLimit)
	}

	r.Run(":8088")
}
