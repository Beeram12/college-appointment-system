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
	Collection         *mongo.Collection
	AvailabilityRecord *mongo.Collection
}

// constructor
func NewAppointment(db *mongo.Database) *MongoAppointment {
	return &MongoAppointment{
		Collection:         db.Collection("appointments"),
		AvailabilityRecord: db.Collection("availability"),
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
	// check if the professor is available at the time slot
	var availabilityRecord models.Availability
	err := m.AvailabilityRecord.FindOne(ctx, bson.M{
		"professor_id": appointment.ProfessorId,
		"time_slot":    appointment.TimeSlot,
		"is_booked":    false,
	}).Decode(&availabilityRecord)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return primitive.NilObjectID, fmt.Errorf("professor is not available at the time slot")
		}
		return primitive.NilObjectID, fmt.Errorf("failed to check professor availability")
	}

	// If the slot is vacant insert into database
	result, err := m.Collection.InsertOne(ctx, appointment)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to book appointment: %w", err)
	}
	appointmentID, good := result.InsertedID.(primitive.ObjectID)
	if !good {
		return primitive.NilObjectID, fmt.Errorf("failed to show insertedID")
	}
	// update the availability record to mark the slot is booked
	updateResult, err := m.AvailabilityRecord.UpdateOne(ctx, bson.M{
		"_id": availabilityRecord.Id,
	}, bson.M{
		"$set": bson.M{"is_booked": true},
	})
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to update professor availability: %w", err)
	}
	if updateResult.ModifiedCount == 0 {
		return primitive.NilObjectID, fmt.Errorf("failed to mark the time slot as booked")
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

// Function for cancelling the appointments
func (m *MongoAppointment) CancelAppointment(ctx context.Context, appointmentID primitive.ObjectID) error {
	var appointment models.Appointment
	// fetch the appointment and verify
	err := m.Collection.FindOne(ctx, bson.M{
		"_id": appointmentID,
	}).Decode(&appointment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("appointments not found")
		}
		return fmt.Errorf("failed to fetch appointment: %w", err)
	}
	// delete the appointment
	result, err := m.Collection.DeleteOne(ctx, bson.M{
		"_id": appointmentID,
	})
	if err != nil {
		return fmt.Errorf("failed to cancel appointment")
	}
	if result.DeletedCount == 0 {
		return errors.New("Appointment not found")
	}
	// update the availability record to mark the slot as unbooked
	_, err = m.AvailabilityRecord.UpdateOne(ctx, bson.M{
		"professor_id": appointment.ProfessorId,
		"time_slot":    appointment.TimeSlot,
	}, bson.M{
		"$set": bson.M{"is_booked": false},
	})
	if err != nil {
		log.Printf("Failed to update the avialability record after the cancellation: %v", err)
	}
	log.Printf("Appointment with %s ID is sucessfully cancelled", appointmentID.Hex())
	return nil
}

// to support ownership verification
func (m *MongoAppointment) GetAppointmentByID(ctx context.Context, appointmentID primitive.ObjectID) (*models.Appointment, error) {
	var appointment models.Appointment
	filter := bson.M{"_id": appointmentID}
	err := m.Collection.FindOne(ctx, filter).Decode(&appointment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("appointment not found")
		}
		return nil, fmt.Errorf("failed to fetch appointment: %w", err)
	}
	return &appointment, nil
}
