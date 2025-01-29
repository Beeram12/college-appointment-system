package repository

import (
	"context"
	"errors"
	"log"

	"github.com/Beeram12/college-appointment-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AppointmentRepo is an interface that defines the methods that should be implemented by the appointment repository
type AppointmentRepo interface {
	BookAppointment(ctx context.Context, appointment models.Appointment) (primitive.ObjectID, error)
	GetAppointmentsByStudentID(ctx context.Context, studentID primitive.ObjectID) ([]models.Appointment, error)
	CancelAppointment(ctx context.Context, appointmentID primitive.ObjectID) error
}

// AppointmentRepository implements the AppointmentRepo interface
type AppointmentRepository struct {
	Collection *mongo.Collection
}

// NewAppointmentRepository creates a new instance of AppointmentRepository
func NewAppointmentRepository(db *mongo.Database) *AppointmentRepository {
	return &AppointmentRepository{
		Collection: db.Collection("appointments"),
	}
}

// BookAppointment inserts a new appointment into the database
func (r *AppointmentRepository) BookAppointment(ctx context.Context, appointment models.Appointment) (primitive.ObjectID, error) {
	result, err := r.Collection.InsertOne(ctx, appointment)
	if err != nil {
		return primitive.NilObjectID, err
	}
	appointmentID := result.InsertedID.(primitive.ObjectID)
	log.Printf("Appointment booked successfully with orderID:%s", appointmentID.Hex())
	return appointmentID, nil
}

// GetAppointmentsByStudentID retrieves all appointments for a given student
func (r *AppointmentRepository) GetAppointmentsByStudentID(ctx context.Context, studentID primitive.ObjectID) ([]models.Appointment, error) {

	var appointments []models.Appointment
	cursor, err := r.Collection.Find(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &appointments); err != nil {
		return nil, err
	}
	return appointments, nil
}

// CancelAppointment deletes an appointment from the database
func (r *AppointmentRepository) CancelAppointment(ctx context.Context, appointmentID primitive.ObjectID) error {

	result, err := r.Collection.DeleteOne(ctx, bson.M{"_id": appointmentID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("appointment not found")
	}
	log.Printf("Appointment with ID %s canceled successfully", appointmentID.Hex())
	return nil
}

// Ensure AppointmentRepository implements AppointmentRepo
var _ AppointmentRepo = (*AppointmentRepository)(nil)
