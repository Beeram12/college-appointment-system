package repository

import (
	"context"

	"github.com/Beeram12/college-appointment-system/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Appointment is an interface for appointment repository

type Appointment interface {
	BookAppointment(ctx context.Context, appointment models.Appointment) (primitive.ObjectID, error)
	GetAppointmentsOfStudent(ctx context.Context, studentID primitive.ObjectID) ([]models.Appointment, error)
	CancelAppointment(ctx context.Context, appointmentID primitive.ObjectID) error
	GetAppointmentByID(ctx context.Context, appointmentID primitive.ObjectID) (*models.Appointment, error)
}
