package router

import (
	"github.com/HEBNUOJ/controller"
	"github.com/HEBNUOJ/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRegisterAndLoginRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CorsMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.POST("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.POST("/api/auth/")
	return r
}

func CollectVerifyRoute(r *gin.Engine) *gin.Engine {
	r1 := r.Group("/api/captcha") // 图形验证码处理路由
	checkCodeController := new(controller.CheckCodeController)
	{
		r1.GET("/refresh", checkCodeController.ReloadVerifyCode)
		r1.GET("/show/:captchaId", checkCodeController.GenVerifyCode)
		r1.GET("/verify", checkCodeController.VerifyCode)
	}
	r2 := r.Group("api/email")
	{
		r2.GET("/refresh", checkCodeController.GenEmailVerifyCode)
		r2.GET("/verify", checkCodeController.VerifyCode)
	}
	return r
}
