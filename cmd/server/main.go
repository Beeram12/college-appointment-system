package main

import (
	"log"
	"net/http"

	"github.com/Beeram12/college-appointment-system/config"
	"github.com/Beeram12/college-appointment-system/internal/handlers"
	"github.com/Beeram12/college-appointment-system/internal/repository"
	"github.com/Beeram12/college-appointment-system/internal/routes"
	"github.com/Beeram12/college-appointment-system/pkg/db"

)

func main() {
	// Load environment variables
	cfg := config.LoadConfig()

	// Initialize MongoDB connection
	client, err := db.InitMongo(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to initialize MongoDB:", err)
	}
	// defer client.Disconnect() // Proper cleanup when exiting

	database := client.Database(cfg.Database)

	// Initialize repositories
	authRepo := repository.NewAuthRepository(database)
	appointmentRepo := repository.NewAppointment(database)
	availabilityRepo := repository.NewAvailability(database)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authRepo)
	appointmentHandler := handlers.NewAppointment(appointmentRepo)
	availabilityHandler := handlers.NewAvailability(availabilityRepo)

	// Setup routes using Gorilla Mux
	router := routes.SetupRoutes(authHandler, appointmentHandler, availabilityHandler)

	// Start the server
	log.Printf("Server running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
