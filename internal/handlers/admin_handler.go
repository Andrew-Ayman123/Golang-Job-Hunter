package handlers
import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Andrew-Ayman123/Job-Hunter/internal/dto"
)

// HandleCreateAdmin handles admin creation (admin-only)
// @Summary Create Admin Account
// @Description Create a new admin account (admin-only endpoint)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param admin body dto.CreateAdminRequest true "Admin creation data"
// @Success 201 {object} dto.CreateUserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/create-admin [post]
func (h *UserHandler) HandleCreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
        h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
        return
    }

	user, err := h.userRepo.CreateAdmin(req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			h.writeErrorResponse(w, "Email already exists", http.StatusConflict)
			return
		}
		h.writeErrorResponse(w, "Failed to create admin: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.CreateUserResponse{
		Message: "Admin account created successfully",
		User:    *user,
	}

	h.writeJSONResponse(w, response, http.StatusCreated)
}

// HandleCreateRecruiter handles recruiter creation (admin-only)
// @Summary Create Recruiter Account
// @Description Create a new recruiter account (admin-only endpoint)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param recruiter body dto.CreateRecruiterRequest true "Recruiter creation data"
// @Success 201 {object} dto.CreateUserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/create-recruiter [post]
func (h *UserHandler) HandleCreateRecruiter(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRecruiterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
        h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
        return
    }

	user, err := h.userRepo.CreateRecruiter(req)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			h.writeErrorResponse(w, "Email already exists", http.StatusConflict)
			return
		}
		h.writeErrorResponse(w, "Failed to create recruiter: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.CreateUserResponse{
		Message: "Recruiter account created successfully",
		User:    *user,
	}

	h.writeJSONResponse(w, response, http.StatusCreated)
}

// HandleCreateCompany handles company creation (admin-only)
// @Summary Create Company
// @Description Create a new company (admin-only endpoint)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param company body dto.CreateCompanyRequest true "Company creation data"
// @Success 201 {object} dto.CreateCompanyResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/create-company [post]
func (h *UserHandler)HandleCreateCompany(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	company, err := h.userRepo.CreateCompany(req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to create company: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.CreateCompanyResponse{
		Message: "Company created successfully",
		Company: *company,
	}

	h.writeJSONResponse(w, response, http.StatusCreated)
}