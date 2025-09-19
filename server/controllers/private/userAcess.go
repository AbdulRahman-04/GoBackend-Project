package private

import (
	"context"
	// "strings"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func UserAccessCollect() {
	userCollection = utils.MongoClient.Database("Event_Booking").Collection("user")
}

// getone user api
func GetOneUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid id format",
		})
		return
	}

	// tokenUserId := c.MustGet("userId").(primitive.ObjectID)
	// if tokenUserId.Hex() != paramId {
	// 	c.JSON(403, gin.H{"msg": "Unauthorized: You can't touch other user's data❌"})
	// 	return
	// }

	// find user by id
	var oneUser models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&oneUser)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Your Profile✨", "OneUser": oneUser})

}

// edit user api
func EditUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userId := c.MustGet("userId").(primitive.ObjectID)

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid id format",
		})
		return
	}

	tokenUserId := c.MustGet("userId").(primitive.ObjectID)
	if tokenUserId.Hex() != paramId {
		c.JSON(403, gin.H{"msg": "Unauthorized: You can't touch other user's data❌"})
		return
	}

	// find user in db
	var editUser models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&editUser)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	// take input
	type User struct {
		UserName string `json:"name" form:"name"`
		// Email    string `json:"email" form:"email"`
		// Password string `json:"password" form:"password"`
		// Phone    string `json:"phone" form:"phone"`
		Language string `json:"language" form:"language"`
		Location string `json:"location" form:"location"`
	}

	// bind to json
	var inputUser User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid request",
		})
		return
	}

	// validations
	if inputUser.UserName == "" || inputUser.Location == "" || inputUser.Language == "" {
		c.JSON(400, gin.H{
			"msg": "Invalid Request, Please add some values to edit ur profile⚠️",
		})
		return
	}

	// update
	update := bson.M{
		"$set": bson.M{
			"username":   inputUser.UserName,
			"language":   inputUser.Language,
			"location":   inputUser.Location,
			"updated_at": time.Now(),
		}}

	// db push
	_, err = userCollection.UpdateByID(ctx, editUser.ID, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Your Profile Updated Successfully!✨",
	})

}

// delete one api
func DeleteOneUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "invalid id format",
		})
		return
	}

	tokenUserId := c.MustGet("userId").(primitive.ObjectID)
	if tokenUserId.Hex() != paramId {
		c.JSON(403, gin.H{"msg": "Unauthorized: You can't touch other user's data❌"})
		return
	}

	// find and delete one
	_, err = userCollection.DeleteOne(ctx, bson.M{"_id": mongoId})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "couldn't delete user, no id found!⚠️",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Your Profile Deleted💔",
	})

}

// User Logout API
func UserLogout(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type LogoutInput struct {
		RefreshToken string `json:"refreshToken"`
	}

	var input LogoutInput
	if err := c.ShouldBindJSON(&input); err != nil || input.RefreshToken == "" {
		c.JSON(400, gin.H{"msg": "Invalid request"})
		return
	}

	// Find user by refresh token
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"refreshToken": input.RefreshToken}).Decode(&user)
	if err != nil {
		c.JSON(401, gin.H{"msg": "Invalid refresh token"})
		return
	}

	// Invalidate refresh token
	_, err = userCollection.UpdateByID(ctx, user.ID, bson.M{
		"$set": bson.M{
			"refreshToken":  "",
			"refreshExpiry": time.Time{},
			"updated_at":    time.Now(),
		},
	})

	if err != nil {
		c.JSON(500, gin.H{"msg": "Could not logout, try again"})
		return
	}

	c.JSON(200, gin.H{
		"msg": "User logged out successfully ✅",
	})
}
