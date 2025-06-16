package handlers
import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/go-chi/chi/v5"

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

	admin, err := h.adminRepo.CreateAdmin(req)
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
		User:    admin.User,
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

	recruiter, err := h.adminRepo.CreateRecruiter(req)
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
		User:    recruiter.User,
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
// @Router /admin/company [post]
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

	company, err := h.adminRepo.CreateCompany(req)
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

// HandleUpdateCompany handles company update (admin-only)
// @Summary Update Company
// @Description Update an existing company (admin-only endpoint)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param companyID path string true "Company ID"
// @Param company body dto.UpdateCompanyRequest true "Company update data"
// @Success 200 {object} dto.CreateCompanyResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/company/{companyID} [patch]
func (h *UserHandler)HandleUpdateCompany(w http.ResponseWriter, r *http.Request){
	var req dto.UpdateCompanyRequest
	companyID := chi.URLParam(r, "companyID")
	if companyID == "" {
		h.writeErrorResponse(w, "Company ID is required", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	//Validation for at least one filed exist and should be at least 6 characters long
	if req.Name == nil && req.Description == nil {
		h.writeErrorResponse(w, "At least one field (name or description) must be provided for update", http.StatusBadRequest)
		return
	}
	if req.Name != nil && len(*req.Name) < 6 {
		h.writeErrorResponse(w, "Company name must be at least 6 characters long", http.StatusBadRequest)
		return
	}
	if req.Description != nil && len(*req.Description) < 6 {
		h.writeErrorResponse(w, "Company description must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	company, err := h.adminRepo.UpdateCompany(companyID, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.writeErrorResponse(w, "Company not found", http.StatusNotFound)
			return
		}
		h.writeErrorResponse(w, "Failed to update company: "+err.Error(), http.StatusInternalServerError)
		return
	}
	response := dto.CreateCompanyResponse{
		Message: "Company updated successfully",
		Company: *company,
	}
	h.writeJSONResponse(w, response, http.StatusOK)

}

// HandleDeleteCompany handles company deletion (admin-only)
// @Summary Delete Company
// @Description Delete an existing company (admin-only endpoint)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param companyID path string true "Company ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/company/{companyID} [delete]
func (h *UserHandler)HandleDeleteCompany(w http.ResponseWriter, r *http.Request){
	// Extract companyID from URL path
	companyID := chi.URLParam(r, "companyID")
	if companyID == "" {
		h.writeErrorResponse(w, "Company ID is required", http.StatusBadRequest)
		return
	}

	err := h.adminRepo.DeleteCompany(companyID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.writeErrorResponse(w, "Company not found", http.StatusNotFound)
			return
		}
		h.writeErrorResponse(w, "Failed to delete company: "+err.Error(), http.StatusInternalServerError)
		return
	}
	h.writeJSONResponse(w, map[string]string{"message": "Company deleted successfully"}, http.StatusOK)
}

