package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Availability struct {
	Id                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProfessorId       primitive.ObjectID `json:"professorId" bson:"professor_id"`
	TimeSlot          time.Time          `json:"timeSlot" bson:"time_slot"`
	TimeSlotFormatted string             `json:"time_slot_formatted" bson:"time_slot_formatted"`
	IsBooked          bool               `json:"isBooked" bson:"is_booked"`
}
