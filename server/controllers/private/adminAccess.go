package private

import (
	"context"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	UserCollection     *mongo.Collection
	EventCollection    *mongo.Collection
	FunctionCollection *mongo.Collection
)

var adminCollection *mongo.Collection

func AdminAccessCollect() {
	adminCollection = utils.MongoClient.Database("Event_Booking").Collection("admin")
}

// ✅ GET ALL USERS
func GetAllUsersAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := UserCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(400, gin.H{"msg": "db decode error"})
		return
	}

	c.JSON(200, gin.H{"msg": "All Users Data", "users": users})
}

// ✅ GET ONE USER
func GetOneUserAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid id format"})
		return
	}

	var oneUser models.User
	err = UserCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&oneUser)
	if err != nil {
		c.JSON(400, gin.H{"msg": "No such user found"})
		return
	}

	c.JSON(200, gin.H{"msg": "One User", "user": oneUser})
}

// ✅ GET ALL EVENTS
func GetAllEventsAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := EventCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}
	defer cursor.Close(ctx)

	var events []models.Event
	if err := cursor.All(ctx, &events); err != nil {
		c.JSON(400, gin.H{"msg": "db decode error"})
		return
	}

	c.JSON(200, gin.H{"msg": "All Events Data", "events": events})
}

// ✅ GET ONE EVENT
func GetOneEventAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid id format"})
		return
	}

	var oneEvent models.Event
	err = EventCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&oneEvent)
	if err != nil {
		c.JSON(400, gin.H{"msg": "No such event found"})
		return
	}

	c.JSON(200, gin.H{"msg": "One Event", "event": oneEvent})
}

// ✅ GET ALL FUNCTIONS
func GetAllFunctionsAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := FunctionCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}
	defer cursor.Close(ctx)

	var functions []models.Function
	if err := cursor.All(ctx, &functions); err != nil {
		c.JSON(400, gin.H{"msg": "db decode error"})
		return
	}

	c.JSON(200, gin.H{"msg": "All Functions Data", "functions": functions})
}

// ✅ GET ONE FUNCTION
func GetOneFunctionAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid id format"})
		return
	}

	var oneFunction models.Function
	err = FunctionCollection.FindOne(ctx, bson.M{"_id": mongoId}).Decode(&oneFunction)
	if err != nil {
		c.JSON(400, gin.H{"msg": "No such function found"})
		return
	}

	c.JSON(200, gin.H{"msg": "One Function", "function": oneFunction})
}

// Admin logout
func AdminLogout(c *gin.Context) {
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

	// Find admin by refresh token
	var admin models.Admin
	err := adminCollection.FindOne(ctx, bson.M{"refreshToken": input.RefreshToken}).Decode(&admin)
	if err != nil {
		c.JSON(401, gin.H{"msg": "Invalid refresh token"})
		return
	}

	// Invalidate refresh token
	_, err = adminCollection.UpdateByID(ctx, admin.ID, bson.M{
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
		"msg": "Admin logged out successfully ✅",
	})
}
