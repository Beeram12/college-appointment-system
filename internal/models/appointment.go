package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ProfessorId primitive.ObjectID `bson:"professor_id"`
	StudentId   primitive.ObjectID `bson:"student_id"`
	TimeSlot    string             `bson:"time_slot"`
}
