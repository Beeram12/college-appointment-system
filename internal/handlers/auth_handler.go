package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Beeram12/college-appointment-system/internal/models"
	"github.com/Beeram12/college-appointment-system/internal/repository"
	"github.com/Beeram12/college-appointment-system/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// auth handler handles authentication
type AuthHandler struct {
	Repo *repository.AuthRepo
}

// NewAuthhandler intializes authhandler

func NewAuthHandler(repo *repository.AuthRepo) *AuthHandler {
	return &AuthHandler{
		Repo: repo,
	}
}

// Register a user
func (m *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}
	// hashing the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword
	user.Id = primitive.NewObjectID()
	//store in db
	err = m.Repo.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered sucessfully"})
}

// Login and generate jwt

func (m *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	// fetching from db
	user, err := m.Repo.FindByUsername(r.Context(), creds.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// check password
	if !utils.CheckPasswordHash(creds.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	// generate jwt token
	token, err := utils.GenerateToken(user.Username, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
