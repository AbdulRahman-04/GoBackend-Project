package middleware

import "github.com/gin-gonic/gin"

// only admins
func OnlyAdmins() gin.HandlerFunc{
	return  func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin"{
			c.JSON(400, gin.H{
				"msg": "Access Denied on this route⚠️, Only Admins allowed here",
			})
			c.Abort()
			return 
		}
		c.Next()
	}
}


// only users
func OnlyUsers() gin.HandlerFunc{
	return  func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "user"{
			c.JSON(400, gin.H{
				"msg": "Access Denied on this route⚠️, Only Users allowed here",
			})
			c.Abort()
			return 
		}
		c.Next()
	}
}