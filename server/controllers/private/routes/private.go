package routes

import (
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/private"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/middleware"
	"github.com/gin-gonic/gin"
)


func PrivateRoutes(r*gin.Engine){
	privateGroup := r.Group("/api/private")
    privateGroup.Use(middleware.AuthMiddleware())
	privateGroup.Use(middleware.RateLimitMiddleware(10))
	 
	{
		// events routes 
		privateGroup.POST("/events/create", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5),private.CreateEvent)
		privateGroup.GET("/getallevents", middleware.OnlyUsers(), middleware.RateLimitMiddleware(10),private.GetAllEvents)
		privateGroup.GET("/getoneevent/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(10),private.GetOneEvent)
		privateGroup.PUT("/updateevent/:id",  middleware.OnlyUsers(), middleware.RateLimitMiddleware(5),private.EditEventApi)
		privateGroup.DELETE("/deleteoneevent/:id",  middleware.OnlyUsers(), middleware.RateLimitMiddleware(5),private.DeleteOneEvent)
		privateGroup.DELETE("/deleteallevents", middleware.OnlyUsers(), middleware.RateLimitMiddleware(1),private.DeleteAllEvents)

		// user logout api
       privateGroup.POST("/users/logout", middleware.OnlyUsers(), private.UserLogout)



		// function routes 
		privateGroup.POST("/func/create", middleware.OnlyUsers(),  middleware.RateLimitMiddleware(5),private.CreateFunction)
		privateGroup.GET("/getallfunc", middleware.OnlyUsers(), middleware.RateLimitMiddleware(10),private.GetAllFunctions)
		privateGroup.GET("/getonefunc/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(10),private.GetOneFunction)
		privateGroup.PUT("/updatefunc/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5),private.EditFunction)
		privateGroup.DELETE("/deleteonefunc/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5),private.DeleteOneFunction)
		privateGroup.DELETE("/deleteallfuncs", middleware.OnlyUsers(), middleware.RateLimitMiddleware(2),private.DeleteAllFunctions)

		// Admins access routes
		privateGroup.GET("/admins/getallevents", middleware.OnlyAdmins(), private.GetAllEventsAdmin)
		privateGroup.GET("/admins/getone/:id", middleware.OnlyAdmins(), private.GetOneEventAdmin)
		privateGroup.GET("/admins/getallusers", middleware.OnlyAdmins(), private.GetAllUsersAdmin)
		privateGroup.GET("/admins/getoneuser/:id", middleware.OnlyAdmins(), private.GetOneUser)
		privateGroup.GET("/admins/getallfuncs", middleware.OnlyAdmins(), private.GetAllFunctionsAdmin)
		privateGroup.GET("/admins/getonefunc/:id", middleware.OnlyAdmins(), private.GetOneFunctionAdmin)
        privateGroup.POST("/admins/logout", middleware.OnlyAdmins(), private.AdminLogout)
	}

}