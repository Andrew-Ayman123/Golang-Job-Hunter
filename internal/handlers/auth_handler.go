package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Andrew-Ayman123/Job-Hunter/internal/dto"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/middleware"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/repository"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo   repository.UserRepository
	jwtService *services.JWTService
}

func NewUserHandler(userRepo repository.UserRepository, jwtService *services.JWTService) *UserHandler {
	return &UserHandler{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// HandleUserSignUp handles user registration
// @Summary User Registration
// @Description Register a new user (applicant or recruiter)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body dto.CreateApplicantRequest true "User registration data"
// @Success 201 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /applicant/signup [post]
func (h *UserHandler) HandleApplicantSignUp(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateApplicantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if err := h.validateApplicantSignUpRequest(req); err != nil {
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.CreateApplicant(req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			h.writeErrorResponse(w, "Email already exists", http.StatusConflict)
			return
		}
		h.writeErrorResponse(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.jwtService.GenerateToken(*user)
	if err != nil {
		h.writeErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := dto.LoginResponse{
		Token: token,
		User:  *user,
	}

	h.writeJSONResponse(w, response, http.StatusCreated)
}

// HandleUserLogIn handles user login
// @Summary User Login
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/login [post]
func (h *UserHandler) HandleUserLogIn(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		h.writeErrorResponse(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		h.writeErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		h.writeErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.jwtService.GenerateToken(*user)
	if err != nil {
		h.writeErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := dto.LoginResponse{
		Token: token,
		User:  *user,
	}

	h.writeJSONResponse(w, response, http.StatusOK)
}

// HandleGetProfile handles getting user profile
// @Summary Get User Profile
// @Description Get the current user's profile
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile [get]
func (h *UserHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		h.writeErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	h.writeJSONResponse(w, user, http.StatusOK)
}

// HandleUpdateProfile handles updating user profile
// @Summary Update User Profile
// @Description Update the current user's profile
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param profile body models.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile [put]
// func (h *UserHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
// 	if !ok {
// 		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	var req models.UpdateProfileRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Get current user
// 	_, err := h.userRepo.GetUserByID(claims.UserID)
// 	if err != nil {
// 		h.writeErrorResponse(w, "User not found", http.StatusNotFound)
// 		return
// 	}

// 	// Update user with the provided data
// 	updatedUser, err := h.userRepo.UpdateUser(claims.UserID, req)
// 	if err != nil {
// 		if strings.Contains(err.Error(), "duplicate key") {
// 			h.writeErrorResponse(w, "Email already exists", http.StatusConflict)
// 			return
// 		}
// 		h.writeErrorResponse(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	h.writeJSONResponse(w, updatedUser, http.StatusOK)
// }

// Helper methods
func (h *UserHandler) validateApplicantSignUpRequest(req dto.CreateApplicantRequest) error {
	if req.Email == "" || req.Password == "" {
		return errors.New("email and password are required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

func (h *UserHandler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
