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

type MongoAvailability struct {
	Collection *mongo.Collection
}

// constructor
func NewAvailability(db *mongo.Database) *MongoAvailability {
	return &MongoAvailability{
		Collection: db.Collection("availability"),
	}
}

// function for adding availability
func (m *MongoAvailability) AddAvailability(ctx context.Context, availability models.Availability) (primitive.ObjectID, error) {
	// check if the availability is already present
	filter := m.Collection.FindOne(ctx, bson.M{
		"professor_id": availability.ProfessorId,
		"time_slot":    availability.TimeSlot,
	})

	if filter.Err() == nil {
		return primitive.NilObjectID, fmt.Errorf("availability already exists for this time slot")
	}

	// Add the slot
	result, err := m.Collection.InsertOne(ctx, availability)
	if err != nil {
		return primitive.NilObjectID, err
	}

	availabilityId, good := result.InsertedID.(primitive.ObjectID)
	if !good {
		return primitive.NilObjectID, fmt.Errorf("failed to show addedId")
	}
	log.Printf("Availability is added sucessfully")
	return availabilityId, nil
}

//function to get all the availabilites of professor

func (m *MongoAvailability) GetAvailabilityOfProfessor(ctx context.Context, professorId primitive.ObjectID) ([]models.Availability, error) {
	var availabilites []models.Availability
	filter := bson.M{"professor_id": professorId, "is_booked": false}

	result, err := m.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer result.Close(ctx)
	// iterating over results and decoding each document into availabilites slice
	for {
		var availability models.Availability
		if !result.Next(ctx) {
			break
		}
		if err := result.Decode(&availability); err != nil {
			return nil, err
		}
		availabilites = append(availabilites, availability)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}
	return availabilites, nil
}

// Deleting availability
func (m *MongoAvailability) DeleteAvailability(ctx context.Context, availabilityId primitive.ObjectID) error {
	result, err := m.Collection.DeleteOne(ctx, bson.M{
		"_id": availabilityId,
	})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Availability not found")
	}
	log.Printf("Availability with %s ID is sucessfully removed", availabilityId.Hex())
	return nil
}
