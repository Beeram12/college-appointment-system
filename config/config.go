package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	Database string
	Port     string
}

type JWT struct {
	JwtSecret string
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
	database := os.Getenv("MONGO_DB")
	if database == "" {
		log.Fatal("Database name is empty")
	}

	// get the port number
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	return &Config{
		MongoURI: mongoURI,
		Database: database,
		Port:     port,
	}
}

func JwtLoadConfig() *JWT {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET")
	return &JWT{
		JwtSecret: jwtSecretKey,
	}
}
