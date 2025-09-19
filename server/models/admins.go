package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Admin struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Role string `bson:"role" json:"role"`

	AdminName string `bson:"adminname" json:"adminname" binding:"required"`
	Email string `bson:"email" json:"email" binding:"required"` // You can optionally add: ,email
	Password string `bson:"password" json:"password" binding:"required,min=6"`
	Phone string  `bson:"phone" json:"phone" binding:"required,min=6"` // Better to use len=10 if fixed length
	Language string `bson:"language" json:"language" binding:"required,oneof=Hindi English Urdu Kannada"`
	Location string `bson:"location" json:"location" binding:"required"`

	AdminVerified struct {
		Email bool `bson:"emailVerified" json:"emailVerified"`
	} `bson:"adminVerified" json:"adminVerified"`

	AdminVerifyToken struct {
		Email string  `bson:"emailVerifyToken" json:"emailVerifyToken"`
		Phone string  `bson:"phoneVerifyToken" json:"phoneVerifyToken"`
	} `bson:"adminVerifyToken" json:"adminVerifyToken"`

	RefreshToken  string    `bson:"refreshToken,omitempty" json:"refreshToken,omitempty"`
	RefreshExpiry time.Time `bson:"refreshExpiry,omitempty" json:"refreshExpiry,omitempty"`

	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
}
