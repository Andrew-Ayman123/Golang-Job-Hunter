// Updated handlers/user.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Andrew-Ayman123/Job-Hunter/internal/dto"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/repository"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/services"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userRepo   repository.UserRepository
	jwtService *services.JWTService
	validator    *validator.Validate
}

func NewUserHandler(userRepo repository.UserRepository, jwtService *services.JWTService) *UserHandler {
	return &UserHandler{
		userRepo:   userRepo,
		jwtService: jwtService,
		validator:  validator.New(),
	}
}

// HandleApplicantSignUp handles applicant registration (public)
// @Summary Applicant Registration
// @Description Register a new applicant (public endpoint)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body dto.CreateApplicantRequest true "Applicant registration data"
// @Success 201 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /applicant/signup [post]
func (h *UserHandler) HandleApplicantSignUp(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateApplicantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
        h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
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