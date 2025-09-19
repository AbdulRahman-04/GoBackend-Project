package routes

import (
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/public"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/middleware"
	"github.com/gin-gonic/gin"
)

func PublicRoutes(r*gin.Engine){
	publicGroup := r.Group("/api/public")
	publicGroup.Use(middleware.RateLimitMiddleware(5))

	{
		// users public apis's
	publicGroup.POST("/users/signup", public.UserSignUp)
	publicGroup.POST("/users/signin", public.UserSignIn)
	publicGroup.GET("/user/emailverify/:token", public.EmailVerifyUser)
	publicGroup.POST("/users/change-password", public.UserChangePass)
	publicGroup.POST("/users/forgot-password", public.UserForgotPass)

	// admins
	publicGroup.POST("/admins/signup", public.AdminSignUp)
	publicGroup.POST("/admins/signin", public.AdminSignIn)
	publicGroup.GET("/admin/emailverify/:token", public.EmailVerifyAdmin)
	publicGroup.POST("/admins/change-password", public.AdminChangePass)
	publicGroup.POST("/admins/forgot-password", public.AdminForgotPass)

	
	}
}