package repository

import (
	"context"

	"github.com/Beeram12/college-appointment-system/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	database *mongo.Database
}

// intializing new repo
func NewAuthRepository(db *mongo.Database) *AuthRepo {
	return &AuthRepo{
		database: db,
	}
}

// Find by username
func (repo *AuthRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	err := repo.database.Collection("users").FindOne(ctx, bson.M{
		"username": username,
	}).Decode(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// If user not present create and insert in the database
func (repo *AuthRepo) CreateUser(ctx context.Context, user models.User) error {
	_, err := repo.database.Collection("users").InsertOne(ctx, user)
	return err
}
