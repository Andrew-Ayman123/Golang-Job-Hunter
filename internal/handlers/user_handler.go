package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/Andrew-Ayman123/Job-Hunter/internal/dto"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/middleware"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/services"
)

// HandleGetProfile handles getting user profile
// @Summary Get User Profile
// @Description Get the current user's profile
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserProfile
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/profile [get]
func (h *UserHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.GetUserProfileByID(claims.UserID)
	if err != nil {
		// write error
		h.writeErrorResponse(w, "Failed to get user profile "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, user, http.StatusOK)
}

// PHONE NUMBERS ENDPOINTS

// @Summary Add Phone Number
// @Description Add a new phone number for the current user
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreatePhoneNumberRequest true "Phone number data"
// @Success 201 {object} models.UserPhoneNumber
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/phone-numbers [post]
func (h *UserHandler) HandleCreatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreatePhoneNumberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	phoneNumber, err := h.userRepo.CreatePhoneNumber(claims.UserID, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to create phone number: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, phoneNumber, http.StatusCreated)
}

// @Summary Update Phone Number
// @Description Update an existing phone number
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param phoneID path string true "Phone Number ID"
// @Param request body dto.CreatePhoneNumberRequest true "Phone number data"
// @Success 200 {object} models.UserPhoneNumber
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/phone-numbers/{phoneID} [put]
func (h *UserHandler) HandleUpdatePhoneNumber(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	phoneID := chi.URLParam(r, "phoneID")
	if phoneID == "" {
		h.writeErrorResponse(w, "Phone number ID is required", http.StatusBadRequest)
		return
	}
	phoneIDProcessed, err := uuid.Parse(phoneID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid phone number ID format", http.StatusBadRequest)
		return
	}

	var req dto.CreatePhoneNumberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	phoneNumber, err := h.userRepo.UpdatePhoneNumber(claims.UserID, phoneIDProcessed, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to update phone number: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, phoneNumber, http.StatusOK)
}

// @Summary Delete Phone Number
// @Description Delete a phone number
// @Tags User Profile
// @Security BearerAuth
// @Param phoneID path string true "Phone Number ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/phone-numbers/{phoneID} [delete]
func (h *UserHandler) HandleDeletePhoneNumber(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	phoneID := chi.URLParam(r, "phoneID")
	if phoneID == "" {
		h.writeErrorResponse(w, "Phone number ID is required", http.StatusBadRequest)
		return
	}

	phoneIDProcessed, err := uuid.Parse(phoneID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid phone number ID format", http.StatusBadRequest)
		return
	}

	err = h.userRepo.DeletePhoneNumber(claims.UserID, phoneIDProcessed)
	if err != nil {
		h.writeErrorResponse(w, "Failed to delete phone number: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// EDUCATION ENDPOINTS

// @Summary Add Education
// @Description Add a new education entry for the current user
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateEducationRequest true "Education data"
// @Success 201 {object} models.UserEducation
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/education [post]
func (h *UserHandler) HandleCreateEducation(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateEducationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	education, err := h.userRepo.CreateEducation(claims.UserID, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to create education: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, education, http.StatusCreated)
}

// @Summary Update Education
// @Description Update an existing education entry
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param educationID path string true "Education ID"
// @Param request body dto.CreateEducationRequest true "Education data"
// @Success 200 {object} models.UserEducation
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/education/{educationID} [put]
func (h *UserHandler) HandleUpdateEducation(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	educationID := chi.URLParam(r, "educationID")
	if educationID == "" {
		h.writeErrorResponse(w, "Education ID is required", http.StatusBadRequest)
		return
	}

	var req dto.CreateEducationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	educationIDProcessed, err := uuid.Parse(educationID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid education ID format", http.StatusBadRequest)
		return
	}

	education, err := h.userRepo.UpdateEducation(claims.UserID, educationIDProcessed, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to update education: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, education, http.StatusOK)
}

// @Summary Delete Education
// @Description Delete an education entry
// @Tags User Profile
// @Security BearerAuth
// @Param educationID path string true "Education ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/education/{educationID} [delete]
func (h *UserHandler) HandleDeleteEducation(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	educationID := chi.URLParam(r, "educationID")
	if educationID == "" {
		h.writeErrorResponse(w, "Education ID is required", http.StatusBadRequest)
		return
	}

	educationIDProcessed, err := uuid.Parse(educationID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid education ID format", http.StatusBadRequest)
		return
	}

	err = h.userRepo.DeleteEducation(claims.UserID, educationIDProcessed)
	if err != nil {
		h.writeErrorResponse(w, "Failed to delete education: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// EXPERIENCE ENDPOINTS

// @Summary Add Experience
// @Description Add a new experience entry for the current user
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateExperienceRequest true "Experience data"
// @Success 201 {object} models.UserExperience
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/experience [post]
func (h *UserHandler) HandleCreateExperience(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateExperienceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	experience, err := h.userRepo.CreateExperience(claims.UserID, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to create experience: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, experience, http.StatusCreated)
}

// @Summary Update Experience
// @Description Update an existing experience entry
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param experienceID path string true "Experience ID"
// @Param request body dto.CreateExperienceRequest true "Experience data"
// @Success 200 {object} models.UserExperience
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/experience/{experienceID} [put]
func (h *UserHandler) HandleUpdateExperience(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	experienceID := chi.URLParam(r, "experienceID")
	if experienceID == "" {
		h.writeErrorResponse(w, "Experience ID is required", http.StatusBadRequest)
		return
	}

	var req dto.CreateExperienceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	experienceIDProcessed, err := uuid.Parse(experienceID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid Experience ID format", http.StatusBadRequest)
		return
	}

	experience, err := h.userRepo.UpdateExperience(claims.UserID, experienceIDProcessed, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to update experience: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, experience, http.StatusOK)
}

// @Summary Delete Experience
// @Description Delete an experience entry
// @Tags User Profile
// @Security BearerAuth
// @Param experienceID path string true "Experience ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/experience/{experienceID} [delete]
func (h *UserHandler) HandleDeleteExperience(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	experienceID := chi.URLParam(r, "experienceID")
	if experienceID == "" {
		h.writeErrorResponse(w, "Experience ID is required", http.StatusBadRequest)
		return
	}

	experienceIDProcessed, err := uuid.Parse(experienceID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid Experience ID format", http.StatusBadRequest)
		return
	}

	err = h.userRepo.DeleteExperience(claims.UserID, experienceIDProcessed)
	if err != nil {
		h.writeErrorResponse(w, "Failed to delete experience: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CERTIFICATION ENDPOINTS

// @Summary Add Certification
// @Description Add a new certification for the current user
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateCertificationRequest true "Certification data"
// @Success 201 {object} models.UserCertification
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/certifications [post]
func (h *UserHandler) HandleCreateCertification(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateCertificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	certification, err := h.userRepo.CreateCertification(claims.UserID, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to create certification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, certification, http.StatusCreated)
}

// @Summary Update Certification
// @Description Update an existing certification
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param certificationID path string true "Certification ID"
// @Param request body dto.CreateCertificationRequest true "Certification data"
// @Success 200 {object} models.UserCertification
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/certifications/{certificationID} [put]
func (h *UserHandler) HandleUpdateCertification(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	certificationID := chi.URLParam(r, "certificationID")
	if certificationID == "" {
		h.writeErrorResponse(w, "Certifcation ID is required", http.StatusBadRequest)
		return
	}

	var req dto.CreateCertificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	certificationIDProcessed, err := uuid.Parse(certificationID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid Certifcation ID format", http.StatusBadRequest)
		return
	}

	certification, err := h.userRepo.UpdateCertification(claims.UserID, certificationIDProcessed, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to update certification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, certification, http.StatusOK)
}

// @Summary Delete Certification
// @Description Delete a certification
// @Tags User Profile
// @Security BearerAuth
// @Param certificationID path string true "Certification ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/certifications/{certificationID} [delete]
func (h *UserHandler) HandleDeleteCertification(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	certificationID := chi.URLParam(r, "certificationID")
	if certificationID == "" {
		h.writeErrorResponse(w, "Certifcation ID is required", http.StatusBadRequest)
		return
	}

	certificationIDProcessed, err := uuid.Parse(certificationID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid Certifcation ID format", http.StatusBadRequest)
		return
	}

	err = h.userRepo.DeleteCertification(claims.UserID, certificationIDProcessed)
	if err != nil {
		h.writeErrorResponse(w, "Failed to delete certification: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PROJECT ENDPOINTS

// @Summary Add Project
// @Description Add a new project for the current user
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateProjectRequest true "Project data"
// @Success 201 {object} models.UserProject
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/projects [post]
func (h *UserHandler) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	project, err := h.userRepo.CreateProject(claims.UserID, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to create project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, project, http.StatusCreated)
}

// @Summary Update Project
// @Description Update an existing project
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param projectID path string true "Project ID"
// @Param request body dto.CreateProjectRequest true "Project data"
// @Success 200 {object} models.UserProject
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/projects/{projectID} [put]
func (h *UserHandler) HandleUpdateProject(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectID := chi.URLParam(r, "projectID")
	if projectID == "" {
		h.writeErrorResponse(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	var req dto.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	projectIDProcessed, err := uuid.Parse(projectID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid Project ID format", http.StatusBadRequest)
		return
	}

	project, err := h.userRepo.UpdateProject(claims.UserID, projectIDProcessed, req)
	if err != nil {
		h.writeErrorResponse(w, "Failed to update project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, project, http.StatusOK)
}

// @Summary Delete Project
// @Description Delete a project
// @Tags User Profile
// @Security BearerAuth
// @Param projectID path string true "Project ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/projects/{projectID} [delete]
func (h *UserHandler) HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectID := chi.URLParam(r, "projectID")
	if projectID == "" {
		h.writeErrorResponse(w, "Project ID is required", http.StatusBadRequest)
		return
	}
	projectIDProcessed, err := uuid.Parse(projectID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid Project ID format", http.StatusBadRequest)
		return
	}

	err = h.userRepo.DeleteProject(claims.UserID, projectIDProcessed)
	if err != nil {
		h.writeErrorResponse(w, "Failed to delete project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SKILLS ENDPOINTS

// @Summary Search Skills
// @Description Search for skills by name
// @Tags Skills
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} models.Skill
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /skills/search [get]
func (h *UserHandler) HandleSearchSkills(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		h.writeErrorResponse(w, "Search query is required", http.StatusBadRequest)
		return
	}

	skills, err := h.userRepo.GetSkillsByName(query)
	if err != nil {
		h.writeErrorResponse(w, "Failed to search skills: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, skills, http.StatusOK)
}

// @Summary Add Skills to User
// @Description Add skills to the current user's profile
// @Tags User Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.AddSkillsRequest true "Skills to add"
// @Success 200 {array} models.Skill
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/skills [post]
func (h *UserHandler) HandleAddUserSkills(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.AddSkillsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if validationErrors, err := h.validateStruct(req); err != nil {
		h.writeJSONResponse(w, validationErrors, http.StatusBadRequest)
		return
	}

	err := h.userRepo.AddUserSkills(claims.UserID, req.SkillIDs)
	if err != nil {
		h.writeErrorResponse(w, "Failed to add skills: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated skills list
	skills, err := h.userRepo.GetUserSkillsByID(claims.UserID)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get updated skills: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSONResponse(w, skills, http.StatusOK)
}

// @Summary Remove Skill from User
// @Description Remove a skill from the current user's profile
// @Tags User Profile
// @Security BearerAuth
// @Param skillID path int true "Skill ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/profile/skills/{skillID} [delete]
func (h *UserHandler) HandleRemoveUserSkill(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
	if !ok {
		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	skillID := chi.URLParam(r, "skillID")
	if skillID == "" {
		h.writeErrorResponse(w, "Skill ID is required", http.StatusBadRequest)
		return
	}

	skillIDProcessed, err := strconv.Atoi(skillID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid Skill ID format", http.StatusBadRequest)
		return
	}

	err = h.userRepo.RemoveUserSkill(claims.UserID, skillIDProcessed)
	if err != nil {
		h.writeErrorResponse(w, "Failed to remove skill: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// // MEDIA UPLOAD ENDPOINTS

// // @Summary Upload Education Media
// // @Description Upload media files for an education entry
// // @Tags Media
// // @Security BearerAuth
// // @Accept multipart/form-data
// // @Produce json
// // @Param id path string true "Education ID"
// // @Param file formData file true "Media file"
// // @Param media_type formData string false "Media type (image, video, document)"
// // @Param alt_text formData string false "Alt text for accessibility"
// // @Param description formData string false "Media description"
// // @Success 201 {object} models.UserMedia
// // @Failure 400 {object} map[string]string
// // @Failure 401 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /user/profile/education/{id}/media [post]
// func (h *UserHandler) HandleUploadEducationMedia(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
// 	if !ok {
// 		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	educationID, err := uuid.Parse(vars["id"])
// 	if err != nil {
// 		h.writeErrorResponse(w, "Invalid education ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Parse multipart form (max 32MB)
// 	err = r.ParseMultipartForm(32 << 20)
// 	if err != nil {
// 		h.writeErrorResponse(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		h.writeErrorResponse(w, "No file provided", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	mediaType := r.FormValue("media_type")
// 	if mediaType == "" {
// 		mediaType = "document" // default
// 	}

// 	altText := r.FormValue("alt_text")
// 	description := r.FormValue("description")

// 	media, err := h.mediaService.UploadEducationMedia(claims.UserID, educationID, file, header.Filename, mediaType, altText, description)
// 	if err != nil {
// 		h.writeErrorResponse(w, "Failed to upload media: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	h.writeJSONResponse(w, media, http.StatusCreated)
// }

// // @Summary Upload Experience Media
// // @Description Upload media files for an experience entry
// // @Tags Media
// // @Security BearerAuth
// // @Accept multipart/form-data
// // @Produce json
// // @Param id path string true "Experience ID"
// // @Param file formData file true "Media file"
// // @Param media_type formData string false "Media type (image, video, document)"
// // @Param alt_text formData string false "Alt text for accessibility"
// // @Param description formData string false "Media description"
// // @Success 201 {object} models.UserMedia
// // @Failure 400 {object} map[string]string
// // @Failure 401 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /user/profile/experience/{id}/media [post]
// func (h *UserHandler) HandleUploadExperienceMedia(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
// 	if !ok {
// 		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	experienceID, err := uuid.Parse(vars["id"])
// 	if err != nil {
// 		h.writeErrorResponse(w, "Invalid experience ID", http.StatusBadRequest)
// 		return
// 	}

// 	err = r.ParseMultipartForm(32 << 20)
// 	if err != nil {
// 		h.writeErrorResponse(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		h.writeErrorResponse(w, "No file provided", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	mediaType := r.FormValue("media_type")
// 	if mediaType == "" {
// 		mediaType = "document"
// 	}

// 	altText := r.FormValue("alt_text")
// 	description := r.FormValue("description")

// 	media, err := h.mediaService.UploadExperienceMedia(claims.UserID, experienceID, file, header.Filename, mediaType, altText, description)
// 	if err != nil {
// 		h.writeErrorResponse(w, "Failed to upload media: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	h.writeJSONResponse(w, media, http.StatusCreated)
// }

// // @Summary Upload Certification Media
// // @Description Upload media files for a certification entry
// // @Tags Media
// // @Security BearerAuth
// // @Accept multipart/form-data
// // @Produce json
// // @Param id path string true "Certification ID"
// // @Param file formData file true "Media file"
// // @Param media_type formData string false "Media type (image, video, document)"
// // @Param alt_text formData string false "Alt text for accessibility"
// // @Param description formData string false "Media description"
// // @Success 201 {object} models.UserMedia
// // @Failure 400 {object} map[string]string
// // @Failure 401 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /user/profile/certifications/{id}/media [post]
// func (h *UserHandler) HandleUploadCertificationMedia(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
// 	if !ok {
// 		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	certificationID, err := uuid.Parse(vars["id"])
// 	if err != nil {
// 		h.writeErrorResponse(w, "Invalid certification ID", http.StatusBadRequest)
// 		return
// 	}

// 	err = r.ParseMultipartForm(32 << 20)
// 	if err != nil {
// 		h.writeErrorResponse(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		h.writeErrorResponse(w, "No file provided", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	mediaType := r.FormValue("media_type")
// 	if mediaType == "" {
// 		mediaType = "document"
// 	}

// 	altText := r.FormValue("alt_text")
// 	description := r.FormValue("description")

// 	media, err := h.mediaService.UploadCertificationMedia(claims.UserID, certificationID, file, header.Filename, mediaType, altText, description)
// 	if err != nil {
// 		h.writeErrorResponse(w, "Failed to upload media: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	h.writeJSONResponse(w, media, http.StatusCreated)
// }

// // @Summary Upload Project Media
// // @Description Upload media files for a project entry
// // @Tags Media
// // @Security BearerAuth
// // @Accept multipart/form-data
// // @Produce json
// // @Param id path string true "Project ID"
// // @Param file formData file true "Media file"
// // @Param media_type formData string false "Media type (image, video, document)"
// // @Param alt_text formData string false "Alt text for accessibility"
// // @Param description formData string false "Media description"
// // @Success 201 {object} models.UserMedia
// // @Failure 400 {object} map[string]string
// // @Failure 401 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /user/profile/projects/{id}/media [post]
// func (h *UserHandler) HandleUploadProjectMedia(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
// 	if !ok {
// 		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	projectID, err := uuid.Parse(vars["id"])
// 	if err != nil {
// 		h.writeErrorResponse(w, "Invalid project ID", http.StatusBadRequest)
// 		return
// 	}

// 	err = r.ParseMultipartForm(32 << 20)
// 	if err != nil {
// 		h.writeErrorResponse(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		h.writeErrorResponse(w, "No file provided", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	mediaType := r.FormValue("media_type")
// 	if mediaType == "" {
// 		mediaType = "document"
// 	}

// 	altText := r.FormValue("alt_text")
// 	description := r.FormValue("description")

// 	media, err := h.mediaService.UploadProjectMedia(claims.UserID, projectID, file, header.Filename, mediaType, altText, description)
// 	if err != nil {
// 		h.writeErrorResponse(w, "Failed to upload media: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	h.writeJSONResponse(w, media, http.StatusCreated)
// }

// // @Summary Delete Media
// // @Description Delete a media file
// // @Tags Media
// // @Security BearerAuth
// // @Param mediaId path string true "Media ID"
// // @Success 204
// // @Failure 400 {object} map[string]string
// // @Failure 401 {object} map[string]string
// // @Failure 404 {object} map[string]string
// // @Failure 500 {object} map[string]string
// // @Router /user/profile/media/{mediaId} [delete]
// func (h *UserHandler) HandleDeleteMedia(w http.ResponseWriter, r *http.Request) {
// 	claims, ok := r.Context().Value(middleware.UserContextKey).(*services.Claims)
// 	if !ok {
// 		h.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	vars := mux.Vars(r)
// 	mediaID, err := uuid.Parse(vars["mediaId"])
// 	if err != nil {
// 		h.writeErrorResponse(w, "Invalid media ID", http.StatusBadRequest)
// 		return
// 	}

// 	err = h.mediaService.DeleteMedia(claims.UserID, mediaID)
// 	if err != nil {
// 		h.writeErrorResponse(w, "Failed to delete media: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }
