package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Availability struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProfessorId primitive.ObjectID `json:"professorId" bson:"professor_id"`
	TimeSlot    []string           `json:"timeSlot" bson:"time_slot"`
}
