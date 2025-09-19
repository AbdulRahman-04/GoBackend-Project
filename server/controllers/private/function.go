package private

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/models"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var functionCollection *mongo.Collection

func FunctionCollect() {
	functionCollection = utils.MongoClient.Database("Event_Booking").Collection("functions")
}

// Create Function
func CreateFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	funcName := c.PostForm("funcname")
	funcType := c.PostForm("functype")
	funcDesc := c.PostForm("funcdes")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	var newFunction models.Function
	newFunction.ID = primitive.NewObjectID()
	newFunction.UserId = userId
	newFunction.FuncName = funcName
	newFunction.FuncType = funcType
	newFunction.FuncDesc = funcDesc
	newFunction.ImageUrl = imageUrl
	newFunction.IsPublic = isPublic
	newFunction.Status = status
	newFunction.Location = location
	newFunction.CreatedAt = time.Now()
	newFunction.UpdatedAt = time.Now()

	_, err = functionCollection.InsertOne(ctx, newFunction)
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"msg": "New Function Created✨", "functionDetails": newFunction})
}


func GetAllFunctions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uid, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"msg": "unauthorized"})
		return
	}

	var userId primitive.ObjectID
	switch v := uid.(type) {
	case string:
		oid, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			c.JSON(500, gin.H{"msg": "invalid userId format"})
			return
		}
		userId = oid
	case primitive.ObjectID:
		userId = v
	default:
		c.JSON(500, gin.H{"msg": "invalid userId type"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit

	cacheKey := fmt.Sprintf("functions:%s:%d:%d", userId.Hex(), page, limit)
	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			var funcs []models.Function
			if jsonErr := json.Unmarshal([]byte(cached), &funcs); jsonErr == nil {
				c.JSON(200, gin.H{
					"msg":       "All Functions are here✨",
					"functions": funcs,
					"page":      page,
					"limit":     limit,
					"source":    "redis",
				})
				fmt.Println("GetAllFunctions: Redis hit")
				return
			}
		}
	}

	// DB fallback
	total, _ := functionCollection.CountDocuments(ctx, bson.M{"userId": userId})
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := functionCollection.Find(ctx, bson.M{"userId": userId}, opts)
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}
	defer cursor.Close(ctx)

	var allFunctions []models.Function
	if err := cursor.All(ctx, &allFunctions); err != nil {
		c.JSON(400, gin.H{"msg": "decoding error"})
		return
	}

	// Cache to Redis
	if utils.RedisClient != nil {
		data, _ := json.Marshal(allFunctions)
		_ = utils.RedisClient.Set(ctx, cacheKey, data, 60*time.Second).Err()
	}

	c.JSON(200, gin.H{
		"msg":       "All Functions are here✨",
		"functions": allFunctions,
		"page":      page,
		"limit":     limit,
		"total":     total,
		"source":    "db",
	})
}

func GetOneFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	uid, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"msg": "unauthorized"})
		return
	}
	userId, ok := uid.(primitive.ObjectID)
	if !ok {
		c.JSON(500, gin.H{"msg": "invalid userId type"})
		return
	}

	cacheKey := fmt.Sprintf("function:%s:%s", userId.Hex(), mongoId.Hex())
	if utils.RedisClient != nil {
		if cached, err := utils.RedisClient.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			var oneFunc models.Function
			if jsonErr := json.Unmarshal([]byte(cached), &oneFunc); jsonErr == nil {
				c.JSON(200, gin.H{
					"msg":      "One Function fetched successfully ✅",
					"function": oneFunc,
					"source":   "redis",
				})
				fmt.Println("GetOneFunction: Redis hit")
				return
			}
		}
	}

	// DB fallback
	var oneFunc models.Function
	err = functionCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oneFunc)
	if err != nil {
		c.JSON(404, gin.H{"msg": "No function found❌"})
		return
	}

	// Cache to Redis
	if utils.RedisClient != nil {
		data, _ := json.Marshal(oneFunc)
		_ = utils.RedisClient.Set(ctx, cacheKey, data, 60*time.Second).Err()
	}

	c.JSON(200, gin.H{
		"msg":      "One Function fetched successfully ✅",
		"function": oneFunc,
		"source":   "db",
	})
	fmt.Println("GetOneFunction: DB fallback & cached")
}


// Edit Function
func EditFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	userId := c.MustGet("userId").(primitive.ObjectID)

	var oldFunc models.Function
	err = functionCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oldFunc)
	if err != nil {
		c.JSON(400, gin.H{"msg": "No function found to update"})
		return
	}

	funcName := c.PostForm("funcname")
	funcType := c.PostForm("functype")
	funcDesc := c.PostForm("funcdes")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	update := bson.M{
		"$set": bson.M{
			"funcname":    funcName,
			"functype":    funcType,
			"funcdes":     funcDesc,
			"ispublic":    isPublic,
			"status":      status,
			"location":    location,
			"imageUrl":    imageUrl,
			"updated_at":  time.Now(),
		}}

	_, err = functionCollection.UpdateByID(ctx, oldFunc.ID, update)
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"msg": "Function Updated Successfully!✅", "updatedFunction": oldFunc})
}

// Delete One Function
func DeleteOneFunction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)
	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param ID"})
		return
	}

	_, err = functionCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId})
	if err != nil {
		c.JSON(400, gin.H{"msg": "No function found to delete"})
		return
	}

	c.JSON(200, gin.H{"msg": "One Function Deleted✅"})
}

// Delete All Functions
func DeleteAllFunctions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.MustGet("userId").(primitive.ObjectID)

	_, err := functionCollection.DeleteMany(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}

	c.JSON(200, gin.H{"msg": "All Functions Deleted✅"})
}
