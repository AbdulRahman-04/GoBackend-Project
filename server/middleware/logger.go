package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		log.Printf("[IVENTS LOG] %v | %3d | %-7s | %s | %s",
			start.Format("2006-01-02 15:04:05"),
			status,
			method,
			path,
			clientIP,
		)

		log.Printf("→ Took %v\n", duration)
	}
}
