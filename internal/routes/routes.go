package routes

import (
	"github.com/Beeram12/college-appointment-system/internal/handlers"
	"github.com/Beeram12/college-appointment-system/internal/middleware"
	"github.com/gorilla/mux"
)

// SetupRoutes initializes the application routes using Gorilla Mux
func SetupRoutes(authHandler *handlers.AuthHandler, appointmentHandler *handlers.AppointmentHandler, availabilityHandler *handlers.AvailabilityHandler) *mux.Router {
	r := mux.NewRouter()

	// 🔹 Public Routes (No authentication required)
	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", authHandler.Register).Methods("POST") // User Registration
	authRoutes.HandleFunc("/login", authHandler.Login).Methods("POST")       // User Login

	// 🔹 Private Routes (Require authentication)
	protectedRoutes := r.PathPrefix("/").Subrouter()
	protectedRoutes.Use(middleware.JWTMiddleware) 

	// Availability Routes
	availabilityRoutes := protectedRoutes.PathPrefix("/availability").Subrouter()
	availabilityRoutes.HandleFunc("/set", availabilityHandler.AddAvailability).Methods("POST")
	availabilityRoutes.HandleFunc("/{professor_id}", availabilityHandler.GetAvailabilityOfProfessor).Methods("GET")

	// Appointment Routes
	appointmentRoutes := protectedRoutes.PathPrefix("/appointments").Subrouter()
	appointmentRoutes.HandleFunc("/book", appointmentHandler.BookAppointment).Methods("POST")
	appointmentRoutes.HandleFunc("/student", appointmentHandler.GetAppointmentsByStudentID).Methods("GET")
	appointmentRoutes.HandleFunc("/{appointment_id}", appointmentHandler.CancelAppointment).Methods("DELETE")

	return r
}
