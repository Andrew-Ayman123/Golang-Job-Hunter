package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	FullName     string    `json:"full_name" db:"full_name"`
	Location     *string   `json:"location,omitempty" db:"location"`
	Role         string    `json:"role" db:"role"` // 'applicant', 'recruiter', or 'admin'
	Title        *string   `json:"title,omitempty" db:"title"`
	AboutSection *string   `json:"about_section,omitempty" db:"about_section"`
	ProfileURL   *string   `json:"profile_url,omitempty" db:"profile_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Company struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Recruiter struct {
	User
	CompanyID *uuid.UUID `json:"company_id,omitempty" db:"company_id"`
}

type Admin struct {
	User
	AdminLevel int       `json:"admin_level" db:"admin_level"`
}

type UserPhoneNumber struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	PhoneType   string    `json:"phone_type" db:"phone_type"` // 'mobile', 'home', 'work', 'other'
	IsPrimary   bool      `json:"is_primary" db:"is_primary"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserEducation struct {
	ID              uuid.UUID   `json:"id" db:"id"`
	UserID          uuid.UUID   `json:"user_id" db:"user_id"`
	InstitutionName string      `json:"institution_name" db:"institution_name"`
	Degree          string      `json:"degree" db:"degree"`
	FieldOfStudy    *string     `json:"field_of_study,omitempty" db:"field_of_study"`
	StartDate       *time.Time  `json:"start_date,omitempty" db:"start_date"`
	EndDate         *time.Time  `json:"end_date,omitempty" db:"end_date"`
	IsCurrent       bool        `json:"is_current" db:"is_current"`
	GradeGPA        *string     `json:"grade_gpa,omitempty" db:"grade_gpa"`
	Description     *string     `json:"description,omitempty" db:"description"`
	Media           []UserMedia `json:"media,omitempty"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
}

type UserExperience struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	UserID         uuid.UUID   `json:"user_id" db:"user_id"`
	CompanyName    string      `json:"company_name" db:"company_name"`
	PositionTitle  string      `json:"position_title" db:"position_title"`
	EmploymentType string      `json:"employment_type" db:"employment_type"` // 'full-time', 'part-time', 'contract', 'internship', 'freelance', 'volunteer'
	StartDate      *time.Time  `json:"start_date,omitempty" db:"start_date"`
	EndDate        *time.Time  `json:"end_date,omitempty" db:"end_date"`
	IsCurrent      bool        `json:"is_current" db:"is_current"`
	Location       *string     `json:"location,omitempty" db:"location"`
	Description    *string     `json:"description,omitempty" db:"description"`
	Media          []UserMedia `json:"media,omitempty"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
}

type UserCertification struct {
	ID                  uuid.UUID   `json:"id" db:"id"`
	UserID              uuid.UUID   `json:"user_id" db:"user_id"`
	CertificationName   string      `json:"certification_name" db:"certification_name"`
	IssuingOrganization string      `json:"issuing_organization" db:"issuing_organization"`
	IssueDate           *time.Time  `json:"issue_date,omitempty" db:"issue_date"`
	ExpirationDate      *time.Time  `json:"expiration_date,omitempty" db:"expiration_date"`
	CredentialID        *string     `json:"credential_id,omitempty" db:"credential_id"`
	CredentialURL       *string     `json:"credential_url,omitempty" db:"credential_url"`
	Description         *string     `json:"description,omitempty" db:"description"`
	Media               []UserMedia `json:"media,omitempty"`
	CreatedAt           time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at" db:"updated_at"`
}

type UserProject struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	UserID           uuid.UUID      `json:"user_id" db:"user_id"`
	ProjectName      string         `json:"project_name" db:"project_name"`
	Description      *string        `json:"description,omitempty" db:"description"`
	StartDate        *time.Time     `json:"start_date,omitempty" db:"start_date"`
	EndDate          *time.Time     `json:"end_date,omitempty" db:"end_date"`
	IsOngoing        bool           `json:"is_ongoing" db:"is_ongoing"`
	ProjectURL       *string        `json:"project_url,omitempty" db:"project_url"`
	Media            []UserMedia    `json:"media,omitempty"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

type UserMedia struct {
	ID              uuid.UUID   `json:"id" db:"id"`
	UserID          uuid.UUID   `json:"user_id" db:"user_id"`
	MediaType       string      `json:"media_type" db:"media_type"` // 'image', 'video', 'document'
	FileName        string      `json:"file_name" db:"file_name"`
	FilePath        string      `json:"file_path" db:"file_path"`
	FileSize        *int64      `json:"file_size,omitempty" db:"file_size"`
	MimeType        *string     `json:"mime_type,omitempty" db:"mime_type"`
	AltText         *string     `json:"alt_text,omitempty" db:"alt_text"`
	Description     *string     `json:"description,omitempty" db:"description"`
	EducationID     *uuid.UUID  `json:"education_id,omitempty" db:"education_id"`
	ExperienceID    *uuid.UUID  `json:"experience_id,omitempty" db:"experience_id"`
	CertificationID *uuid.UUID  `json:"certification_id,omitempty" db:"certification_id"`
	ProjectID       *uuid.UUID  `json:"project_id,omitempty" db:"project_id"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
}

type Skill struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
}

// Composite structs for API responses with related data
type UserProfile struct {
	User           User                `json:"user"`
	PhoneNumbers   []UserPhoneNumber   `json:"phone_numbers,omitempty"`
	Education      []UserEducation     `json:"education,omitempty"`
	Experience     []UserExperience    `json:"experience,omitempty"`
	Certifications []UserCertification `json:"certifications,omitempty"`
	Projects       []UserProject       `json:"projects,omitempty"`
	Skills         []Skill             `json:"skills,omitempty"`
}
