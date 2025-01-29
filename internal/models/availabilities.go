package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Availability struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	ProfessorId primitive.ObjectID `bson:"professor_id"`
	TimeSlot    []string           `bson:"time_slot"`
}
