package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name     string             `json:"name,omitempty" bson:"name"`
	Email    string             `json:"email,omitempty" bson:"email"`
	Password string             `json:"password,omitempty" bson:"password"`
	Role     string             `json:"role,omitempty" bson:"role"`
}
