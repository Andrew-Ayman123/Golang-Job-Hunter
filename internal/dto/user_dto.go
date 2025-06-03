package dto

import (
	models "github.com/Andrew-Ayman123/Job-Hunter/internal/models"
)

type CreateApplicantRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	FullName *string `json:"full_name,omitempty"`
	// Role     string `json:"role" validate:"required,oneof=applicant recruiter"`

}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}
