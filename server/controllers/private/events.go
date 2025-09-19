package private

import (
	"context"
	"encoding/json"
	"fmt"

	// "encoding/json"
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

var eventsCollection *mongo.Collection

func EventsCollect() {
	eventsCollection = utils.MongoClient.Database("Event_Booking").Collection("events")
}

// create even api
func CreateEvent(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get userId
	userId := c.MustGet("userId").(primitive.ObjectID)

	// take input from form
	eventName := c.PostForm("eventname")
	eventType := c.PostForm("eventtype")
	eventAttendenceStr := c.PostForm("attendence")
	eventAttendenceInt, err := strconv.Atoi(eventAttendenceStr)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Conversion error",
		})
		return
	}
	eventDes := c.PostForm("eventdesc")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	// var bnake push in db
	var newEvent models.Event

	newEvent.ID = primitive.NewObjectID()
	newEvent.UserId = userId
	newEvent.EventName = eventName
	newEvent.EventtType = eventType
	newEvent.EventAttendence = eventAttendenceInt
	newEvent.EventDescription = eventDes
	newEvent.IsPublic = isPublic
	newEvent.Status = status
	newEvent.Location = location
	newEvent.ImageUrl = imageUrl
	newEvent.CreatedAt = time.Now()
	newEvent.UpdatedAt = time.Now()

	_, err = eventsCollection.InsertOne(ctx, newEvent)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "New Event Created✨", "event Details": newEvent})

}

// Get All Events with Pagination + Redis + Source info
func GetAllEvents(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ---------------- Safe userId handling ----------------
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

	// ---------------- Pagination params ----------------
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit

	// ---------------- Redis cache check ----------------
	cacheKey := fmt.Sprintf("events:%s:%d:%d", userId.Hex(), page, limit)
	cachedData, err := utils.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedResponse struct {
			Msg     string         `json:"msg"`
			Events  []models.Event `json:"events"`
			Page    int            `json:"page"`
			Limit   int            `json:"limit"`
			Total   int64          `json:"total"`
			HasNext bool           `json:"hasNext"`
			HasPrev bool           `json:"hasPrev"`
			Source  string         `json:"source"`
		}
		if jsonErr := json.Unmarshal([]byte(cachedData), &cachedResponse); jsonErr == nil {
			cachedResponse.Source = "redis"
			fmt.Println("GetAllEvents: Data served from Redis")
			c.JSON(200, cachedResponse)
			return
		}
	}

	// ---------------- DB fallback ----------------
	total, err := eventsCollection.CountDocuments(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(500, gin.H{"msg": "failed to count events"})
		return
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cursor, err := eventsCollection.Find(ctx, bson.M{"userId": userId}, opts)
	if err != nil {
		c.JSON(400, gin.H{"msg": "db error"})
		return
	}
	defer cursor.Close(ctx)

	var allEvents []models.Event
	if err = cursor.All(ctx, &allEvents); err != nil {
		c.JSON(400, gin.H{"msg": "decoding error"})
		return
	}

	// ---------------- Prepare response ----------------
	response := struct {
		Msg     string         `json:"msg"`
		Events  []models.Event `json:"events"`
		Page    int            `json:"page"`
		Limit   int            `json:"limit"`
		Total   int64          `json:"total"`
		HasNext bool           `json:"hasNext"`
		HasPrev bool           `json:"hasPrev"`
		Source  string         `json:"source"`
	}{
		Msg:     "All Events Are here✨",
		Events:  allEvents,
		Page:    page,
		Limit:   limit,
		Total:   total,
		HasNext: int64(skip+limit) < total,
		HasPrev: page > 1,
		Source:  "db",
	}

	// ---------------- Save to Redis (without source field) ----------------
	cacheResp := response
	cacheResp.Source = "" // remove before caching
	dataBytes, _ := json.Marshal(cacheResp)
	_ = utils.RedisClient.Set(ctx, cacheKey, dataBytes, 60*time.Second).Err() // TTL 60s

	fmt.Println("GetAllEvents: Data served from DB and cached in Redis")
	c.JSON(200, response)
}

// Get One Event with Redis + source info
func GetOneEvent(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ---------------- Param ID ----------------
	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{"msg": "Invalid param Id"})
		return
	}

	// ---------------- User ID ----------------
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

	cacheKey := fmt.Sprintf("event:%s:%s", userId.Hex(), mongoId.Hex())

	// ---------------- Redis cache check ----------------
	if utils.RedisClient != nil {
		cachedData, err := utils.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil && cachedData != "" {
			var cachedEvent models.Event
			if jsonErr := json.Unmarshal([]byte(cachedData), &cachedEvent); jsonErr == nil {
				c.JSON(200, gin.H{
					"msg":      "One event fetched successfully ✅",
					"OneEvent": cachedEvent,
					"source":   "redis",
				})
				fmt.Println("GetOneEvent: Data served from Redis")
				return
			}
		}
	}

	// ---------------- DB fallback ----------------
	var oneEvent models.Event
	err = eventsCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&oneEvent)
	if err != nil {
		c.JSON(404, gin.H{"msg": "No event found ❌"})
		return
	}

	// ---------------- Save to Redis ----------------
	if utils.RedisClient != nil {
		dataBytes, _ := json.Marshal(oneEvent)
		_ = utils.RedisClient.Set(ctx, cacheKey, dataBytes, 60*time.Second).Err()
	}

	// ---------------- Response ----------------
	c.JSON(200, gin.H{
		"msg":      "One event fetched successfully ✅",
		"OneEvent": oneEvent,
		"source":   "db",
	})
	fmt.Println("GetOneEvent: Data served from DB and cached in Redis")
}


// edit event api
func EditEventApi(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid param Id",
		})
		return
	}

	userId := c.MustGet("userId").(primitive.ObjectID)

	// find it in db
	var editEvent models.Event
	err = eventsCollection.FindOne(ctx, bson.M{"userId": userId, "_id": mongoId}).Decode(&editEvent)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid db error",
		})
		return
	}

	// take inputs

	eventName := c.PostForm("eventname")
	eventType := c.PostForm("eventtype")
	eventAttendenceStr := c.PostForm("attendence")
	eventAttendenceInt, err := strconv.Atoi(eventAttendenceStr)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Conversion error",
		})
		return
	}
	eventDes := c.PostForm("eventdesc")
	isPublic := c.PostForm("ispublic")
	status := c.PostForm("status")
	location := c.PostForm("location")
	imageUrl, err := utils.FileUpload(c)
	if err != nil {
		imageUrl = ""
	}

	// var bnake push in db
	// var myEvent models.Event

	//  update db
	update := bson.M{
		"$set": bson.M{
			"eventname":  eventName,
			"eventtype":  eventType,
			"attendence": eventAttendenceInt,
			"eventdesc":  eventDes,
			"ispublic":   isPublic,
			"status":     status,
			"location":   location,
			"imageUrl":   imageUrl,
			"updated_at": time.Now(),
		}}
	// update the db
	_, err = eventsCollection.UpdateByID(ctx, mongoId, update)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "db error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "Event Updated Successfully!✅", "UpdatedEvent": editEvent})
}

// DeleteOne Event Api
func DeleteOneEvent(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userId
	userId := c.MustGet("userId").(primitive.ObjectID)

	paramId := c.Param("id")
	mongoId, err := primitive.ObjectIDFromHex(paramId)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "Invalid param Id",
		})
		return
	}

	// find one id and delete
	// var deleteEvent models.Event
	_, err = eventsCollection.DeleteOne(ctx, bson.M{"userId": userId, "_id": mongoId})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "No Event Found or userid not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "One Event is deleted✅",
	})
}

// delete all events apis
func DeleteAllEvents(c *gin.Context) {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// userid
	userId := c.MustGet("userId").(primitive.ObjectID)

	// find userId and delete all events of it
	_, err := eventsCollection.DeleteMany(ctx, bson.M{"userId": userId})
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "DB error",
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "All Events Deleted✅",
	})
}
