package handlers

import (
	"net/http"

	"github.com/Andrew-Ayman123/Job-Hunter/internal/middleware"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/services"
	
)

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