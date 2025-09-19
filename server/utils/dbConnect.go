package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient*mongo.Client

func DBConnect() error {
	// ctx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connection request send
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.DBURI))
	if err != nil {
		fmt.Println("Coundlnt connect to mongoDb", err)
	}

	// send test request to mongodb 
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println("Test Request failed", err)
	}

	MongoClient = client
	fmt.Println("MongoDB Connected Successfully!✅")
	return  nil
}