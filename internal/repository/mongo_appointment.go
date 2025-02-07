package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Beeram12/college-appointment-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoAppointment struct {
	Collection *mongo.Collection
}

// constructor
func NewAppointment(db *mongo.Database) *MongoAppointment {
	return &MongoAppointment{
		Collection: db.Collection("appointments"),
	}
}

// Function for booking appointments
func (m *MongoAppointment) BookAppointment(ctx context.Context, appointment models.Appointment) (primitive.ObjectID, error) {
	// check if the professor has appointment already on the slot
	existingAppointment := m.Collection.FindOne(ctx, bson.M{
		"professor_id": appointment.ProfessorId,
		"time_slot":    appointment.TimeSlot,
	})
	
	if existingAppointment.Err() == nil {
		return primitive.NilObjectID, fmt.Errorf("professor already has an appointment on this slot")
	}
	// If the slot is vacant
	result, err := m.Collection.InsertOne(ctx, appointment)
	if err != nil {
		return primitive.NilObjectID, err
	}
	appointmentID, good := result.InsertedID.(primitive.ObjectID)
	if !good {
		return primitive.NilObjectID, fmt.Errorf("failed to show insertedID")
	}
	log.Printf("Appointment is booked sucessfully")
	return appointmentID, nil
}

// Function for getting all the appointments of the student
func (m *MongoAppointment) GetAppointmentsOfStudent(ctx context.Context, studentID primitive.ObjectID) ([]models.Appointment, error) {
	var appointments []models.Appointment

	// filter to find appointments with studentId
	filter := bson.M{"student_id": studentID}

	result, err := m.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer result.Close(ctx)
	// iterating over results and decodes each document into appointments slice
	for {
		var appointment models.Appointment
		if !result.Next(ctx) {
			break
		}
		if err := result.Decode(&appointment); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

// Function for getting all the appointments of professor
func (m *MongoAppointment) GetAppointmentsOfProfessor(ctx context.Context, professorID primitive.ObjectID) ([]models.Appointment, error) {
	var appointments []models.Appointment

	// filter to find appointments with studentId
	filter := bson.M{"professor_id": professorID}

	result, err := m.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer result.Close(ctx)
	// iterating over results and decodes each document into appointments slice
	for {
		var appointment models.Appointment
		if !result.Next(ctx) {
			break
		}
		if err := result.Decode(&appointment); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}
	return appointments, nil
}

// Function for cancelling the appointments
func (m *MongoAppointment) CancelAppointment(ctx context.Context, appointmentID primitive.ObjectID) error {
	result, err := m.Collection.DeleteOne(ctx, bson.M{
		"_id": appointmentID,
	})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Appointment not found")
	}
	log.Printf("Appointment with %s ID is sucessfully cancelled", appointmentID.Hex())
	return nil
}
