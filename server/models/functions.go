package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Function struct{
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId primitive.ObjectID `bson:"userId" json:"userId"`
	FuncName string `bson:"funcname" json:"funcname" binding:"required,min=5,max=20"`
	FuncType string `bson:"functype" json:"functype" binding:"required,oneof=Shaadi Valima Sanchak BabyShower Manjay Aqeeqa"`
	FuncDesc string `bson:"funcdes" json:"funcdes" binding:"required,min=10,max=100"`
	ImageUrl string `bson:"imageUrl" json:"imageUrl" binding:"required"`
	IsPublic string `bson:"ispublic" json:"ispublic" binding:"required,oneof=public private"`
	Status string `bson:"status" json:"status" binding:"required,oneof=Upcoming Cancelled Completed"`
	Location string `bson:"location" json:"location" binding:"required,min=15,max=100"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}