package models

import (
	"time"
	"github.com/google/uuid"
)	

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"` // 'applicant' or 'recruiter'
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Applicant struct {
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	FullName        string    `json:"full_name" db:"full_name"`
	Resume          *string   `json:"resume,omitempty" db:"resume"`
	ExperienceYears *int      `json:"experience_years,omitempty" db:"experience_years"`
	Location        *string   `json:"location,omitempty" db:"location"`
}

type Recruiter struct {
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	CompanyName   string    `json:"company_name" db:"company_name"`
	ContactNumber *string   `json:"contact_number,omitempty" db:"contact_number"`
	Position      *string   `json:"position,omitempty" db:"position"`
}

