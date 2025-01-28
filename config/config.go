package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	Database string
}

func LoadConfig() *Config {
	// Loaded the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the value of mongoURI from .env file
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("URI string of mongoDB is empty")
	}

	// get the value of database name from .env file
	database := os.Getenv("DATABASE")
	if database == "" {
		log.Fatal("Database name is empty")
	}

	return &Config{
		MongoURI: mongoURI,
		Database: database,
	}
}
