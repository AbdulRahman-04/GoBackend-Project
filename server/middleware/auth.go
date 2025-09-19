package middleware

import (
	"strings"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var myKey = []byte(config.AppConfig.JWTKEY)

func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
        
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(400, gin.H{
				"msg": "No token provided",
			})
			c.Abort()
			return 
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer"{
			c.JSON(400, gin.H{
				"msg": "invalid token format",
			})
			c.Abort()
			return 
		}

		myToken := parts[1]

		token, err := jwt.Parse(myToken, func(t *jwt.Token) (interface{}, error) {
			return myKey, nil
		})
		if err != nil {
			c.JSON(400, gin.H{
				"msg": "Invalid or Expired Token❌",
			})
			c.Abort()
			return 
		}

		// get data from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(400, gin.H{
				"msg": "No Data found in token",
			})
			c.Abort()
			return 
		}

		// get userId from token 
		userStrId, ok := claims["id"].(string)
		if !ok {
			c.JSON(400, gin.H{
				"msg": "No userId Data found in token",
			})
			c.Abort()
			return 
		}

		userId, err := primitive.ObjectIDFromHex(userStrId)
		if err != nil  {
			c.JSON(400, gin.H{
				"msg": "error converting to user id or exxpired token",
			})
			c.Abort()
			return 
		}

		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(400, gin.H{
				"msg": "No role Data found in token",
			})
			c.Abort()
			return 
		}

		// set the role and userid in context variable
		c.Set("userId", userId)
		c.Set("role", role)

		c.Next()
	}
}