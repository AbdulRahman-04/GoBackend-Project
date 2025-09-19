package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Role string `bson:"role" json:"role"`

	Username string `bson:"username" json:"username" binding:"required"`
	Email string `bson:"email" json:"email" binding:"required"` // Optionally add: ,email
	Password string `bson:"password" json:"password" binding:"required,min=6"`
	Phone string  `bson:"phone" json:"phone" binding:"required,min=6"` // Or use len=10
	Language string `bson:"language" json:"language" binding:"required,oneof=Hindi English Urdu Kannada"`
	Location string `bson:"location" json:"location" binding:"required"`

	Userverified struct {
		Email bool `bson:"emailVerified" json:"emailVerified"`
	} `bson:"userverified" json:"userverified"`

	Userverifytoken struct {
		Email string  `bson:"emailVerifyToken" json:"emailVerifyToken"`
		Phone string  `bson:"phoneVerifyToken" json:"phoneVerifyToken"`
	} `bson:"userverifytoken" json:"userverifytoken"`

	// -------------- ADD THESE FOR REFRESH TOKEN ----------------
	RefreshToken  string    `bson:"refreshToken" json:"refreshToken"`
	RefreshExpiry time.Time `bson:"refreshExpiry" json:"refreshExpiry"`

	Createdat time.Time  `bson:"created_at" json:"created_at"`
	Updatedat time.Time  `bson:"updated_at" json:"updated_at"`
}
