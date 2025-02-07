package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Beeram12/college-appointment-system/internal/middleware"
	"github.com/Beeram12/college-appointment-system/internal/models"
	"github.com/Beeram12/college-appointment-system/internal/repository"
	"github.com/Beeram12/college-appointment-system/pkg/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AvailabilityHandler handles the availability-related requests
type AvailabilityHandler struct {
	Repo repository.Availability
}

// Initializes a new availability handler
func NewAvailability(repo repository.Availability) *AvailabilityHandler {
	return &AvailabilityHandler{
		Repo: repo,
	}
}

// Utility function to extract user claims from request context
func getUserClaims(r *http.Request) (*utils.CustomClaims, bool) {
	userClaims, ok := r.Context().Value(middleware.UserContextKey).(*utils.CustomClaims)
	return userClaims, ok
}

// Utility function to send error responses
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}

// AddAvailability handles adding availability for a professor
func (m *AvailabilityHandler) AddAvailability(w http.ResponseWriter, r *http.Request) {
	claims, ok := getUserClaims(r)
	if !ok || claims.Role != "professor" {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	professorID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid professor ID")
		return
	}
	var req struct {
		TimeSlot string `json:"time_slot"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	availability := models.Availability{
		ProfessorId: professorID,
		TimeSlot:    req.TimeSlot,
		IsBooked:    false,
	}
	availabilityID, err := m.Repo.AddAvailability(r.Context(), availability)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to add availability")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"availability_id": availabilityID.Hex(),
	})
}

// GetAvailability handles fetching the availability for a specific professor
func (m *AvailabilityHandler) GetAvailabilityOfProfessor(w http.ResponseWriter, r *http.Request) {
	// Extract professor ID from URL path
	vars := mux.Vars(r)
	professorIDHex := strings.TrimSpace(vars["professor_id"])
	professorID, err := primitive.ObjectIDFromHex(professorIDHex)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid professor ID")
		return
	}
	log.Printf("Converted professor ID: %v", professorID)
	// Fetch availability from the repository
	availabilities, err := m.Repo.GetAvailabilityOfProfessor(r.Context(), professorID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch availability")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(availabilities)
}

// DeleteAvailability handles removing an availability for a professor
func (m *AvailabilityHandler) DeleteAvailability(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := getUserClaims(r)
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Ensuring the role is professor
	if userClaims.Role != "professor" {
		sendErrorResponse(w, http.StatusForbidden, "Permission denied")
		return
	}

	// Extract availability ID from URL path
	vars := mux.Vars(r)
	availabilityID, err := primitive.ObjectIDFromHex(vars["availability_id"])
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "Invalid availability ID")
		return
	}

	// Perform delete operation
	err = m.Repo.DeleteAvailability(r.Context(), availabilityID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to delete availability")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Availability deleted successfully"))
}
