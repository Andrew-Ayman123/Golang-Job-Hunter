// Updated dto/user.go
package dto

import (
	models "github.com/Andrew-Ayman123/Job-Hunter/internal/models"
	"github.com/google/uuid"
)

type CreateApplicantRequest struct {
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	FullName *string `json:"full_name,omitempty"`
}

type CreateAdminRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	AdminLevel int `json:"admin_level" validate:"required,min=1,max=5"`
}

type CreateRecruiterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	CompanyID     *uuid.UUID `json:"company_id" validate:"required"`
}
type CreateCompanyRequest struct {
	Name       string `json:"name" validate:"required"`
	Description    string `json:"description" validate:"required"`
}
type CreateCompanyResponse struct {
	Message string `json:"message"`
	Company models.Company `json:"company"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type CreateUserResponse struct {
	Message string      `json:"message"`
	User    models.User `json:"user"`
}