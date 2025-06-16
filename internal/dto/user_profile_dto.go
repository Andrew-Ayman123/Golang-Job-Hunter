package dto

// DTO structures for requests
type CreatePhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	PhoneType   string `json:"phone_type" validate:"required,oneof=mobile home work other"`
	IsPrimary   bool   `json:"is_primary" validate:"required"`
}

type CreateEducationRequest struct {
	InstitutionName string `json:"institution_name" validate:"required"`
	Degree          string `json:"degree" validate:"required"`
	FieldOfStudy    string `json:"field_of_study"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	IsCurrent       bool   `json:"is_current"`
	GradeGPA        string `json:"grade_gpa"`
	Description     string `json:"description"`
}


type CreateExperienceRequest struct {
	CompanyName    string `json:"company_name" validate:"required"`
	PositionTitle  string `json:"position_title" validate:"required"`
	EmploymentType string `json:"employment_type" validate:"required,oneof=full-time part-time contract internship freelance volunteer"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	IsCurrent      bool   `json:"is_current"`
	Location       string `json:"location"`
	Description    string `json:"description"`
}

type CreateCertificationRequest struct {
	CertificationName   string `json:"certification_name" validate:"required"`
	IssuingOrganization string `json:"issuing_organization" validate:"required"`
	IssueDate           string `json:"issue_date"`
	ExpirationDate      string `json:"expiration_date"`
	CredentialID        string `json:"credential_id"`
	CredentialURL       string `json:"credential_url"`
	Description         string `json:"description"`
}

type CreateProjectRequest struct {
	ProjectName string `json:"project_name" validate:"required"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	IsOngoing   bool   `json:"is_ongoing"`
	ProjectURL  string `json:"project_url"`
}

type AddSkillsRequest struct {
	SkillIDs []int `json:"skill_ids" validate:"required,min=1"`
}