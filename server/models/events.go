package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId           primitive.ObjectID `bson:"userId" json:"userId"`

	EventName        string `bson:"eventname" json:"eventname" binding:"required,min=5,max=30"`
	EventtType       string `bson:"eventtype" json:"eventtype" binding:"required,oneof=Party Bar Birthday Gettogether Formal"`
	EventAttendence  int    `bson:"attendence" json:"attendence" binding:"required,min=1"` // removed 'numeric', int already
	EventDescription string `bson:"eventdesc" json:"eventdesc" binding:"required"`
	ImageUrl         string `bson:"imageUrl" json:"imageUrl" binding:"required"`

	IsPublic string `bson:"ispublic" json:"ispublic" binding:"required,oneof=public private"`
	Status string `bson:"status" json:"status" binding:"required,oneof=Upcoming Cancelled Completed"`
	Location string `bson:"location" json:"location" binding:"required,min=15,max=100"`

	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at" json:"updated_at"`
}
