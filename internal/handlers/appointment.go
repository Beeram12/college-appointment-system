package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

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
	userClaims, ok := r.Context().Value(middleware.UserContextKey).(*utils.CustomClaims)
	if !ok || userClaims.Role != "student" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	studentID, err := primitive.ObjectIDFromHex(userClaims.UserID)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	var req struct {
		ProfessorId string `json:"professor_id"`
		TimeSlot    string `json:"time_sloot"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	professorID, err := primitive.ObjectIDFromHex(req.ProfessorId)
	if err != nil {
		http.Error(w, "Invalid professor ID", http.StatusBadRequest)
		return
	}
	parsedTime, err := time.Parse(utils.TimeLayout, req.TimeSlot)
	if err != nil {
		http.Error(w, "Invalid time format.Use 'hh:mm AM/PM'", http.StatusBadRequest)
		return
	}
	appointment := models.Appointment{
		ProfessorId: professorID,
		StudentId:   studentID,
		TimeSlot:    parsedTime,
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
		"time_slot":      parsedTime.Format(utils.TimeLayout),
	})
}

func (m *AppointmentHandler) CancelAppointment(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value(middleware.UserContextKey).(*utils.CustomClaims)
	if !ok || userClaims.Role != "professor" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	appointmentIDHex := strings.TrimSpace(vars["appointment_id"])
	appointmentID, err := primitive.ObjectIDFromHex(appointmentIDHex)
	if err != nil {
		http.Error(w, "Invalid appointmentID", http.StatusBadRequest)
		return
	}
	// fetch the appointment to verify ownership
	appointment, err := m.Repo.GetAppointmentByID(r.Context(), appointmentID)
	if err != nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}
	professorID, err := primitive.ObjectIDFromHex(userClaims.UserID)
	if err != nil {
		http.Error(w, "Invalid Professor ID", http.StatusBadRequest)
		return
	}
	// Ensure the professor is authorized to cancel the appointment
	if appointment.ProfessorId != professorID {
		http.Error(w, "You are not authorized to cancel this appointment", http.StatusForbidden)
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

func (m *AppointmentHandler) GetAppointmentsByStudentID(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := r.Context().Value(middleware.UserContextKey).(*utils.CustomClaims)
	// ensuring the role is student
	if !ok || userClaims.Role != "student" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Convert student ID from jwt claims

	studentID, err := primitive.ObjectIDFromHex(userClaims.UserID)
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
	for i := range appointments {
		appointments[i].TimeSlotFormatted = appointments[i].TimeSlot.Format(utils.TimeLayout)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}
