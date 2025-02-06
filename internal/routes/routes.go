package routes

import (
	"net/http"

	"github.com/Beeram12/college-appointment-system/internal/handlers"
	"github.com/Beeram12/college-appointment-system/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(appointmentHandler *handlers.AppointmentHandler, availabilityHandler *handlers.AvailabilityHandler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware for authentication
	r.Use(middleware.JWTMiddleware)

	// Availability Routes
	r.Route("/availability", func(r chi.Router) {
		r.Post("/set", availabilityHandler.AddAvailability)
		r.Get("/{professor_id}", availabilityHandler.GetAvailabilityOfProfessor)
	})

	// Appointment Routes
	r.Route("/appointments", func(r chi.Router) {
		r.Post("/book", appointmentHandler.BookAppointment)
		r.Get("/student", appointmentHandler.GetAppointmentsByStudentID)
		r.Delete("/{appointment_id}", appointmentHandler.CancelAppointment)
	})

	// Health Check Route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running!"))
	})

	return r
}
