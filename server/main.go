package main

import (
	"fmt"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/private"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/private/routes"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/public"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/middleware"
	// "github.com/AbdulRahman-04/GoProjects/EventManagement/server/routes"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	// ----------------- DB + Redis -----------------
	utils.DBConnect()
	utils.ConnectRedis()

	router := gin.Default()

	// ----------------- Custom Logger -----------------
	router.Use(middleware.CustomLogger())

	// ----------------- Secure Headers (Helmet equivalent) -----------------
	router.Use(SecureHeaders())

	// ----------------- Rate Limiter -----------------
	store := memory.NewStore()
	rate, _ := limiter.NewRateFromFormatted("100-S") // 100 requests/sec
	instance := limiter.New(store, rate)
	router.Use(ginlimiter.NewMiddleware(instance))

	
	// ----------------- CORS -----------------
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://yourdomain.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ----------------- Function calls -----------------
	public.UserCollect()
	public.AdminCollect()
	private.UserAccessCollect()
	private.EventsCollect()
	private.FunctionCollect()
	private.AdminAccessCollect()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "Hello World From Gin"})
	})

	// ----------------- Routes register -----------------
	routes.PublicRoutes(router)
	routes.PrivateRoutes(router)

	// ----------------- Run server -----------------
	router.Run(fmt.Sprintf(":%d", config.AppConfig.Port))
}

// ----------------- Secure Headers Middleware -----------------
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		c.Next();
	}
}
