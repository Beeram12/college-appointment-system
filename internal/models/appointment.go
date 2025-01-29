package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Professor primitive.ObjectID `json:"professorId,omitempty" bson:"professorId,omitempty"`
	Student   primitive.ObjectID `json:"studentId,omitempty" bson:"studentId,omitempty"`
	TimeSlot  string             `json:"timeSlot,omitempty" bson:"timeSlot,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt int64              `bson:"created_at,omitempty"`
}
