// Updated repository/user.go
package repository

import (
	"database/sql"
	"fmt"

	"github.com/Andrew-Ayman123/Job-Hunter/internal/dto"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(req dto.CreateUserRequest, role string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserProfileByID(id uuid.UUID) (*models.UserProfile, error)
	GetUserPhoneNumbersByID(userID uuid.UUID) ([]models.UserPhoneNumber, error)
	GetUserEducationByID(userID uuid.UUID) ([]models.UserEducation, error)
	GetUserExperienceByID(userID uuid.UUID) ([]models.UserExperience, error)
	GetUserCertificationsByID(userID uuid.UUID) ([]models.UserCertification, error)
	GetUserProjectsByID(userID uuid.UUID) ([]models.UserProject, error)
	GetUserSkillsByID(userID uuid.UUID) ([]models.Skill, error)
	GetSkillsByName(name string) ([]models.Skill, error)

	// Phone Numbers
	CreatePhoneNumber(userID uuid.UUID, req dto.CreatePhoneNumberRequest) (*models.UserPhoneNumber, error)
	UpdatePhoneNumber(userID, phoneID uuid.UUID, req dto.CreatePhoneNumberRequest) (*models.UserPhoneNumber, error)
	DeletePhoneNumber(userID, phoneID uuid.UUID) error

	// Education
	CreateEducation(userID uuid.UUID, req dto.CreateEducationRequest) (*models.UserEducation, error)
	UpdateEducation(userID, educationID uuid.UUID, req dto.CreateEducationRequest) (*models.UserEducation, error)
	DeleteEducation(userID, educationID uuid.UUID) error

	// Experience
	CreateExperience(userID uuid.UUID, req dto.CreateExperienceRequest) (*models.UserExperience, error)
	UpdateExperience(userID, experienceID uuid.UUID, req dto.CreateExperienceRequest) (*models.UserExperience, error)
	DeleteExperience(userID, experienceID uuid.UUID) error

	// Certifications
	CreateCertification(userID uuid.UUID, req dto.CreateCertificationRequest) (*models.UserCertification, error)
	UpdateCertification(userID, certificationID uuid.UUID, req dto.CreateCertificationRequest) (*models.UserCertification, error)
	DeleteCertification(userID, certificationID uuid.UUID) error

	// Projects
	CreateProject(userID uuid.UUID, req dto.CreateProjectRequest) (*models.UserProject, error)
	UpdateProject(userID, projectID uuid.UUID, req dto.CreateProjectRequest) (*models.UserProject, error)
	DeleteProject(userID, projectID uuid.UUID) error

	// Skills
	AddUserSkills(userID uuid.UUID, skillIDs []int) error
	RemoveUserSkill(userID uuid.UUID, skillID int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(req dto.CreateUserRequest, role string) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create user
	var user models.User
	query := `
		INSERT INTO users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, password_hash, full_name, role, created_at, updated_at
	`
	err = tx.QueryRow(query, req.Email, string(hashedPassword), req.FullName, role).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if req.FullName == nil {
		return nil, fmt.Errorf("full_name is required for applicants")
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, full_name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, full_name, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserProfileByID(id uuid.UUID) (*models.UserProfile, error) {
	// Get the base user information
	user, err := r.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	profile := &models.UserProfile{
		User: *user,
	}

	// Get all related data concurrently or sequentially
	phoneNumbers, err := r.GetUserPhoneNumbersByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get phone numbers: %w", err)
	}
	profile.PhoneNumbers = phoneNumbers

	education, err := r.GetUserEducationByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get education: %w", err)
	}
	profile.Education = education

	experience, err := r.GetUserExperienceByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get experience: %w", err)
	}
	profile.Experience = experience

	certifications, err := r.GetUserCertificationsByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get certifications: %w", err)
	}
	profile.Certifications = certifications

	projects, err := r.GetUserProjectsByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	profile.Projects = projects

	skills, err := r.GetUserSkillsByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get skills: %w", err)
	}
	profile.Skills = skills

	return profile, nil
}

func (r *userRepository) GetUserPhoneNumbersByID(userID uuid.UUID) ([]models.UserPhoneNumber, error) {
	query := `
		SELECT id, user_id, phone_number, phone_type, is_primary, created_at, updated_at
		FROM user_phone_numbers
		WHERE user_id = $1
		ORDER BY is_primary DESC, created_at ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phoneNumbers []models.UserPhoneNumber
	for rows.Next() {
		var phone models.UserPhoneNumber
		err := rows.Scan(
			&phone.ID, &phone.UserID, &phone.PhoneNumber, &phone.PhoneType,
			&phone.IsPrimary, &phone.CreatedAt, &phone.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		phoneNumbers = append(phoneNumbers, phone)
	}

	return phoneNumbers, rows.Err()
}

func (r *userRepository) GetUserEducationByID(userID uuid.UUID) ([]models.UserEducation, error) {
	query := `
		SELECT id, user_id, institution_name, degree, field_of_study, start_date, end_date,
			   is_current, grade_gpa, description, created_at, updated_at
		FROM user_education
		WHERE user_id = $1
		ORDER BY is_current DESC, end_date DESC NULLS FIRST, start_date DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var education []models.UserEducation
	for rows.Next() {
		var edu models.UserEducation
		err := rows.Scan(
			&edu.ID, &edu.UserID, &edu.InstitutionName, &edu.Degree, &edu.FieldOfStudy,
			&edu.StartDate, &edu.EndDate, &edu.IsCurrent, &edu.GradeGPA, &edu.Description,
			&edu.CreatedAt, &edu.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get media for this education entry
		media, err := r.getMediaForEducation(edu.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get media for education %s: %w", edu.ID, err)
		}
		edu.Media = media

		education = append(education, edu)
	}

	return education, rows.Err()
}

func (r *userRepository) GetUserExperienceByID(userID uuid.UUID) ([]models.UserExperience, error) {
	query := `
		SELECT id, user_id, company_name, position_title, employment_type, start_date, end_date,
			   is_current, location, description, created_at, updated_at
		FROM user_experience
		WHERE user_id = $1
		ORDER BY is_current DESC, end_date DESC NULLS FIRST, start_date DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var experience []models.UserExperience
	for rows.Next() {
		var exp models.UserExperience
		err := rows.Scan(
			&exp.ID, &exp.UserID, &exp.CompanyName, &exp.PositionTitle, &exp.EmploymentType,
			&exp.StartDate, &exp.EndDate, &exp.IsCurrent, &exp.Location, &exp.Description,
			&exp.CreatedAt, &exp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get media for this experience entry
		media, err := r.getMediaForExperience(exp.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get media for experience %s: %w", exp.ID, err)
		}
		exp.Media = media

		experience = append(experience, exp)
	}

	return experience, rows.Err()
}

func (r *userRepository) GetUserCertificationsByID(userID uuid.UUID) ([]models.UserCertification, error) {
	query := `
		SELECT id, user_id, certification_name, issuing_organization, issue_date, expiration_date,
			   credential_id, credential_url, description, created_at, updated_at
		FROM user_certifications
		WHERE user_id = $1
		ORDER BY issue_date DESC NULLS LAST, created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var certifications []models.UserCertification
	for rows.Next() {
		var cert models.UserCertification
		err := rows.Scan(
			&cert.ID, &cert.UserID, &cert.CertificationName, &cert.IssuingOrganization,
			&cert.IssueDate, &cert.ExpirationDate, &cert.CredentialID, &cert.CredentialURL,
			&cert.Description, &cert.CreatedAt, &cert.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get media for this certification entry
		media, err := r.getMediaForCertification(cert.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get media for certification %s: %w", cert.ID, err)
		}
		cert.Media = media

		certifications = append(certifications, cert)
	}

	return certifications, rows.Err()
}

func (r *userRepository) GetUserProjectsByID(userID uuid.UUID) ([]models.UserProject, error) {
	query := `
		SELECT id, user_id, project_name, description, start_date, end_date, is_ongoing,
			   project_url, created_at, updated_at
		FROM user_projects
		WHERE user_id = $1
		ORDER BY is_ongoing DESC, end_date DESC NULLS FIRST, start_date DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.UserProject
	for rows.Next() {
		var project models.UserProject
		err := rows.Scan(
			&project.ID, &project.UserID, &project.ProjectName, &project.Description,
			&project.StartDate, &project.EndDate, &project.IsOngoing, &project.ProjectURL, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Get media for this project entry
		media, err := r.getMediaForProject(project.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get media for project %s: %w", project.ID, err)
		}
		project.Media = media

		projects = append(projects, project)
	}

	return projects, rows.Err()
}

func (r *userRepository) GetUserSkillsByID(userID uuid.UUID) ([]models.Skill, error) {
	query := `
		SELECT s.id, s.name
		FROM skills s
		INNER JOIN user_skills us ON s.id = us.skill_id
		WHERE us.user_id = $1
		ORDER BY s.name ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user_skills []models.Skill
	for rows.Next() {
		var user_skill models.Skill
		err := rows.Scan(&user_skill.ID, &user_skill.Name)
		if err != nil {
			return nil, err
		}
		user_skills = append(user_skills, user_skill)
	}

	return user_skills, rows.Err()
}

// GetSkillsByName queries the database for user skills matching the given name.
func (r *userRepository) GetSkillsByName(name string) ([]models.Skill, error) {
	var skills []models.Skill

	query := `SELECT id, name FROM skills WHERE name LIKE '%$1%' ORDER BY name ASC`
	rows, err := r.db.Query(query, name)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var skill models.Skill
		if err := rows.Scan(&skill.ID, &skill.Name); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return skills, nil
}

// Helper functions to get media for specific entities
func (r *userRepository) getMediaForEducation(educationID uuid.UUID) ([]models.UserMedia, error) {
	return r.getMedia("education_id", educationID)
}

func (r *userRepository) getMediaForExperience(experienceID uuid.UUID) ([]models.UserMedia, error) {
	return r.getMedia("experience_id", experienceID)
}

func (r *userRepository) getMediaForCertification(certificationID uuid.UUID) ([]models.UserMedia, error) {
	return r.getMedia("certification_id", certificationID)
}

func (r *userRepository) getMediaForProject(projectID uuid.UUID) ([]models.UserMedia, error) {
	return r.getMedia("project_id", projectID)
}

func (r *userRepository) getMedia(fieldName string, entityID uuid.UUID) ([]models.UserMedia, error) {
	query := fmt.Sprintf(`
		SELECT id, user_id, media_type, file_name, file_path, file_size, mime_type,
			   alt_text, description, education_id, experience_id, certification_id, project_id,
			   created_at, updated_at
		FROM user_media
		WHERE %s = $1
		ORDER BY created_at ASC
	`, fieldName)

	rows, err := r.db.Query(query, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var media []models.UserMedia
	for rows.Next() {
		var m models.UserMedia
		err := rows.Scan(
			&m.ID, &m.UserID, &m.MediaType, &m.FileName, &m.FilePath, &m.FileSize,
			&m.MimeType, &m.AltText, &m.Description, &m.EducationID, &m.ExperienceID,
			&m.CertificationID, &m.ProjectID, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		media = append(media, m)
	}

	return media, rows.Err()
}

func (r *userRepository) CreatePhoneNumber(userID uuid.UUID, req dto.CreatePhoneNumberRequest) (*models.UserPhoneNumber, error) {
	query := `
  INSERT INTO user_phone_numbers (user_id, phone_number, phone_type, is_primary)
  VALUES ($1, $2, $3, $4)
  RETURNING id, user_id, phone_number, phone_type, is_primary, created_at, updated_at
 `
	var phone models.UserPhoneNumber
	err := r.db.QueryRow(query, userID, req.PhoneNumber, req.PhoneType, req.IsPrimary).Scan(
		&phone.ID, &phone.UserID, &phone.PhoneNumber, &phone.PhoneType,
		&phone.IsPrimary, &phone.CreatedAt, &phone.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create phone number: %w", err)
	}
	return &phone, nil
}

func (r *userRepository) UpdatePhoneNumber(userID, phoneID uuid.UUID, req dto.CreatePhoneNumberRequest) (*models.UserPhoneNumber, error) {
	query := `
  UPDATE user_phone_numbers
  SET phone_number = $1, phone_type = $2, is_primary = $3, updated_at = NOW()
  WHERE user_id = $4 AND id = $5
  RETURNING id, user_id, phone_number, phone_type, is_primary, created_at, updated_at
 `
	var phone models.UserPhoneNumber
	err := r.db.QueryRow(query, req.PhoneNumber, req.PhoneType, req.IsPrimary, userID, phoneID).Scan(
		&phone.ID, &phone.UserID, &phone.PhoneNumber, &phone.PhoneType,
		&phone.IsPrimary, &phone.CreatedAt, &phone.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update phone number: %w", err)
	}
	return &phone, nil
}
func (r *userRepository) DeletePhoneNumber(userID, phoneID uuid.UUID) error {
	query := `
  DELETE FROM user_phone_numbers
  WHERE user_id = $1 AND id = $2
 `
	result, err := r.db.Exec(query, userID, phoneID)
	if err != nil {
		return fmt.Errorf("failed to delete phone number: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (r *userRepository) CreateEducation(userID uuid.UUID, req dto.CreateEducationRequest) (*models.UserEducation, error) {
	query := `
  INSERT INTO user_education (user_id, institution_name, degree, field_of_study,
	  start_date, end_date, is_current, grade_gpa, description)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
  RETURNING id, user_id, institution_name, degree, field_of_study,
	  start_date, end_date, is_current, grade_gpa, description,
	  created_at, updated_at
 `
	var edu models.UserEducation
	err := r.db.QueryRow(query,
		userID, req.InstitutionName, req.Degree, req.FieldOfStudy,
		req.StartDate, req.EndDate, req.IsCurrent, req.GradeGPA,
		req.Description).Scan(
		&edu.ID, &edu.UserID, &edu.InstitutionName, &edu.Degree,
		&edu.FieldOfStudy, &edu.StartDate, &edu.EndDate,
		&edu.IsCurrent, &edu.GradeGPA, &edu.Description,
		&edu.CreatedAt, &edu.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create education: %w", err)
	}
	return &edu, nil
}
func (r *userRepository) UpdateEducation(userID, educationID uuid.UUID, req dto.CreateEducationRequest) (*models.UserEducation, error) {
	query := `
  UPDATE user_education
  SET institution_name = $1, degree = $2, field_of_study = $3,
	  start_date = $4, end_date = $5, is_current = $6,
	  grade_gpa = $7, description = $8, updated_at = NOW()
  WHERE user_id = $9 AND id = $10
  RETURNING id, user_id, institution_name, degree, field_of_study,
	  start_date, end_date, is_current, grade_gpa, description,
	  created_at, updated_at
 `
	var edu models.UserEducation
	err := r.db.QueryRow(query,
		req.InstitutionName, req.Degree, req.FieldOfStudy,
		req.StartDate, req.EndDate, req.IsCurrent, req.GradeGPA,
		req.Description, userID, educationID).Scan(
		&edu.ID, &edu.UserID, &edu.InstitutionName, &edu.Degree,
		&edu.FieldOfStudy, &edu.StartDate, &edu.EndDate,
		&edu.IsCurrent, &edu.GradeGPA, &edu.Description,
		&edu.CreatedAt, &edu.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update education: %w", err)
	}
	return &edu, nil
}
func (r *userRepository) DeleteEducation(userID, educationID uuid.UUID) error {
	query := `
  DELETE FROM user_education
  WHERE user_id = $1 AND id = $2
 `
	result, err := r.db.Exec(query, userID, educationID)
	if err != nil {
		return fmt.Errorf("failed to delete education: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (r *userRepository) CreateExperience(userID uuid.UUID, req dto.CreateExperienceRequest) (*models.UserExperience, error) {
	query := `
  INSERT INTO user_experience (user_id, company_name, position_title, employment_type,
   start_date, end_date, is_current, location, description)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
  RETURNING id, user_id, company_name, position_title, employment_type,
   start_date, end_date, is_current, location, description,
   created_at, updated_at
 `
	var exp models.UserExperience
	err := r.db.QueryRow(query,
		userID, req.CompanyName, req.PositionTitle, req.EmploymentType,
		req.StartDate, req.EndDate, req.IsCurrent, req.Location,
		req.Description).Scan(
		&exp.ID, &exp.UserID, &exp.CompanyName, &exp.PositionTitle,
		&exp.EmploymentType, &exp.StartDate, &exp.EndDate,
		&exp.IsCurrent, &exp.Location, &exp.Description,
		&exp.CreatedAt, &exp.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create experience: %w", err)
	}
	return &exp, nil
}
func (r *userRepository) UpdateExperience(userID, experienceID uuid.UUID, req dto.CreateExperienceRequest) (*models.UserExperience, error) {
	query := `
  UPDATE user_experience
  SET company_name = $1, position_title = $2, employment_type = $3,
   start_date = $4, end_date = $5, is_current = $6,
   location = $7, description = $8, updated_at = NOW()
  WHERE user_id = $9 AND id = $10
  RETURNING id, user_id, company_name, position_title, employment_type,
   start_date, end_date, is_current, location, description,
   created_at, updated_at
 `
	var exp models.UserExperience
	err := r.db.QueryRow(query,
		req.CompanyName, req.PositionTitle, req.EmploymentType,
		req.StartDate, req.EndDate, req.IsCurrent, req.Location,
		req.Description, userID, experienceID).Scan(
		&exp.ID, &exp.UserID, &exp.CompanyName, &exp.PositionTitle,
		&exp.EmploymentType, &exp.StartDate, &exp.EndDate,
		&exp.IsCurrent, &exp.Location, &exp.Description,
		&exp.CreatedAt, &exp.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update experience: %w", err)
	}
	return &exp, nil
}
func (r *userRepository) DeleteExperience(userID, experienceID uuid.UUID) error {
	query := `
  DELETE FROM user_experience
  WHERE user_id = $1 AND id = $2
 `
	result, err := r.db.Exec(query, userID, experienceID)
	if err != nil {
		return fmt.Errorf("failed to delete experience: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (r *userRepository) CreateCertification(userID uuid.UUID, req dto.CreateCertificationRequest) (*models.UserCertification, error) {
	query := `
  INSERT INTO user_certifications (user_id, certification_name, issuing_organization,
   issue_date, expiration_date, credential_id, credential_url, description)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
  RETURNING id, user_id, certification_name, issuing_organization,
   issue_date, expiration_date, credential_id, credential_url,
   description, created_at, updated_at
 `
	var cert models.UserCertification
	err := r.db.QueryRow(query,
		userID, req.CertificationName, req.IssuingOrganization,
		req.IssueDate, req.ExpirationDate, req.CredentialID,
		req.CredentialURL, req.Description).Scan(
		&cert.ID, &cert.UserID, &cert.CertificationName,
		&cert.IssuingOrganization, &cert.IssueDate,
		&cert.ExpirationDate, &cert.CredentialID,
		&cert.CredentialURL, &cert.Description,
		&cert.CreatedAt, &cert.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create certification: %w", err)
	}
	return &cert, nil
}
func (r *userRepository) UpdateCertification(userID, certificationID uuid.UUID, req dto.CreateCertificationRequest) (*models.UserCertification, error) {
	query := `
  UPDATE user_certifications
  SET certification_name = $1, issuing_organization = $2,
   issue_date = $3, expiration_date = $4, credential_id = $5,
   credential_url = $6, description = $7, updated_at = NOW()
  WHERE user_id = $8 AND id = $9
  RETURNING id, user_id, certification_name, issuing_organization,
   issue_date, expiration_date, credential_id, credential_url,
   description, created_at, updated_at
 `
	var cert models.UserCertification
	err := r.db.QueryRow(query,
		req.CertificationName, req.IssuingOrganization,
		req.IssueDate, req.ExpirationDate, req.CredentialID,
		req.CredentialURL, req.Description, userID, certificationID).Scan(
		&cert.ID, &cert.UserID, &cert.CertificationName,
		&cert.IssuingOrganization, &cert.IssueDate,
		&cert.ExpirationDate, &cert.CredentialID,
		&cert.CredentialURL, &cert.Description,
		&cert.CreatedAt, &cert.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update certification: %w", err)
	}
	return &cert, nil
}
func (r *userRepository) DeleteCertification(userID, certificationID uuid.UUID) error {
	query := `
  DELETE FROM user_certifications
  WHERE user_id = $1 AND id = $2
 `
	result, err := r.db.Exec(query, userID, certificationID)
	if err != nil {
		return fmt.Errorf("failed to delete certification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (r *userRepository) CreateProject(userID uuid.UUID, req dto.CreateProjectRequest) (*models.UserProject, error) {
	query := `
  INSERT INTO user_projects (user_id, project_name, description, start_date,
   end_date, is_ongoing, project_url)
  VALUES ($1, $2, $3, $4, $5, $6, $7)
  RETURNING id, user_id, project_name, description,
   start_date, end_date, is_ongoing, project_url,
   created_at, updated_at
 `
	var project models.UserProject
	err := r.db.QueryRow(query,
		userID, req.ProjectName, req.Description,
		req.StartDate, req.EndDate, req.IsOngoing,
		req.ProjectURL).Scan(
		&project.ID, &project.UserID, &project.ProjectName,
		&project.Description, &project.StartDate,
		&project.EndDate, &project.IsOngoing,
		&project.ProjectURL, &project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return &project, nil
}
func (r *userRepository) UpdateProject(userID, projectID uuid.UUID, req dto.CreateProjectRequest) (*models.UserProject, error) {
	query := `
  UPDATE user_projects
  SET project_name = $1, description = $2, start_date = $3,
   end_date = $4, is_ongoing = $5, project_url = $6,
   updated_at = NOW()
  WHERE user_id = $7 AND id = $8
  RETURNING id, user_id, project_name, description,
   start_date, end_date, is_ongoing, project_url,
   created_at, updated_at
 `
	var project models.UserProject
	err := r.db.QueryRow(query,
		req.ProjectName, req.Description, req.StartDate,
		req.EndDate, req.IsOngoing, req.ProjectURL,
		userID, projectID).Scan(
		&project.ID, &project.UserID, &project.ProjectName,
		&project.Description, &project.StartDate,
		&project.EndDate, &project.IsOngoing,
		&project.ProjectURL, &project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	return &project, nil
}
func (r *userRepository) DeleteProject(userID, projectID uuid.UUID) error {
	query := `
  DELETE FROM user_projects
  WHERE user_id = $1 AND id = $2
 `
	result, err := r.db.Exec(query, userID, projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
func (r *userRepository) AddUserSkills(userID uuid.UUID, skillIDs []int) error {
	query := `
  INSERT INTO user_skills (user_id, skill_id)
  VALUES ($1, $2)
  ON CONFLICT (user_id, skill_id) DO NOTHING
 `
	for _, skillID := range skillIDs {
		_, err := r.db.Exec(query, userID, skillID)
		if err != nil {
			return fmt.Errorf("failed to add skill %d for user %s: %w", skillID, userID, err)
		}
	}
	return nil
}
func (r *userRepository) RemoveUserSkill(userID uuid.UUID, skillID int) error {
	query := `
  DELETE FROM user_skills
  WHERE user_id = $1 AND skill_id = $2
 `
	result, err := r.db.Exec(query, userID, skillID)
	if err != nil {
		return fmt.Errorf("failed to remove skill %d for user %s: %w", skillID, userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
