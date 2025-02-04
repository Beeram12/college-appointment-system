package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Appointment struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProfessorId primitive.ObjectID `json:"professorId" bson:"professor_id"`
	StudentId   primitive.ObjectID `json:"studentId" bson:"student_id"`
	TimeSlot    string             `json:"timeslot" bson:"time_slot"`
}
