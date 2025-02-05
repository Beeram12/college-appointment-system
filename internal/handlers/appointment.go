package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Beeram12/college-appointment-system/internal/middleware"
	"github.com/Beeram12/college-appointment-system/internal/models"
	"github.com/Beeram12/college-appointment-system/internal/repository"
	"github.com/Beeram12/college-appointment-system/pkg/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// handler handles the appointment related requests
type AppointmentHandler struct {
	Repo repository.Appointment
}

// intilizaes a new appointment handler
func NewAppointment(repo repository.Appointment) *AppointmentHandler {
	return &AppointmentHandler{
		Repo: repo,
	}
}

func (m *AppointmentHandler) BookAppointment(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserContextKey).(*utils.CustomClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var appointment models.Appointment
	err := json.NewDecoder(r.Body).Decode(&appointment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// insert apponitment into db
	appointmentID, err := m.Repo.BookAppointment(r.Context(), appointment)
	if err != nil {
		http.Error(w, "Failed to Book appointment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"appointment_id": appointmentID.Hex(),
	})
}

func (m *AppointmentHandler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value(middleware.UserContextKey).(*utils.CustomClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	appointmentID, err := primitive.ObjectIDFromHex(vars["appointment_id"])
	if err != nil {
		http.Error(w, "Invalid appointmentID", http.StatusBadRequest)
		return
	}
	// check roles based access
	if userClaims.Role != "professor" {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	// perform cancel operation
	err = m.Repo.CancelAppointment(r.Context(), appointmentID)
	if err != nil {
		http.Error(w, "Failed to cancel appointment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Appointment canceled successfully"))
}

func (m *AppointmentHandler) GetAppointmentsByStudentID(w http.ResponseWriter, r http.Request) {
	userClaims, ok := r.Context().Value(middleware.UserContextKey).(*utils.CustomClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ensuring the role is student
	if userClaims.Role != "student" {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	// Convert student ID from jwt claims
	studentID, err := primitive.ObjectIDFromHex(userClaims.Username)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusUnauthorized)
		return
	}
	// Fetching
	appointments, err := m.Repo.GetAppointmentsOfStudent(r.Context(), studentID)
	if err != nil {
		http.Error(w, "Failed to fetch appointments", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}


