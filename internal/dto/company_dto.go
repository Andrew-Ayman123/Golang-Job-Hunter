package dto

import (
	models "github.com/Andrew-Ayman123/Job-Hunter/internal/models"
)

type CreateCompanyRequest struct {
	Name        string `json:"name" validate:"required,min=6"`
	Description string `json:"description" validate:"required,min=6"`
}
type CreateCompanyResponse struct {
	Message string         `json:"message"`
	Company models.Company `json:"company"`
}

type UpdateCompanyRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}