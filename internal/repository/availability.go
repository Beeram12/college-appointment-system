package repository

import (
	"context"

	"github.com/Beeram12/college-appointment-system/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Availability interface{
	AddAvailability(ctx context.Context,availability models.Availability)(primitive.ObjectID,error)
	GetAvailabilityOfProfessor(ctx context.Context,professorId primitive.ObjectID)([]models.Availability,error)
	DeleteAvailability(ctx context.Context,availabilityId primitive.ObjectID)error
}
