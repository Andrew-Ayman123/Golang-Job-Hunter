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
	CreateApplicant(req dto.CreateApplicantRequest) (*models.User, error)
	CreateAdmin(req dto.CreateAdminRequest) (*models.User, error)
	CreateRecruiter(req dto.CreateRecruiterRequest) (*models.User, error)
	CreateCompany(req dto.CreateCompanyRequest) (*models.Company, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateApplicant(req dto.CreateApplicantRequest) (*models.User, error) {
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
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, 'applicant')
		RETURNING id, email, password_hash, role, created_at, updated_at
	`
	err = tx.QueryRow(query, req.Email, string(hashedPassword)).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if req.FullName == nil {
		return nil, fmt.Errorf("full_name is required for applicants")
	}
	applicant := models.Applicant{
		UserID: user.ID,
		Resume: nil, // Resume can be added later
	}
	err = r.createApplicantTx(tx, applicant)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user, nil
}

func (r *userRepository) CreateAdmin(req dto.CreateAdminRequest) (*models.User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with admin role
	var user models.User
	query := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, 'admin')
		RETURNING id, email, password_hash, role, created_at, updated_at
	`
	err = r.db.QueryRow(query, req.Email, string(hashedPassword)).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) CreateRecruiter(req dto.CreateRecruiterRequest) (*models.User, error) {
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

	// Create user with recruiter role
	var user models.User
	query := `
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, 'recruiter')
		RETURNING id, email, password_hash, role, created_at, updated_at
	`
	err = tx.QueryRow(query, req.Email, string(hashedPassword)).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create recruiter user: %w", err)
	}

	// Create recruiter profile
	recruiter := models.Recruiter{
		UserID:      user.ID,
		CompanyID: req.CompanyID,
	}
	err = r.createRecruiterTx(tx, recruiter)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateCompany(req dto.CreateCompanyRequest) (*models.Company, error) {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	// Create company
	var company models.Company
	query := `
		INSERT INTO companies (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at
	`
	err = tx.QueryRow(query, req.Name, req.Description).Scan(
		&company.ID, &company.Name, &company.Description, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &company, nil
}

func (r *userRepository) createApplicantTx(tx *sql.Tx, applicant models.Applicant) error {
	query := `
		INSERT INTO applicants (user_id)
		VALUES ($1)
	`
	_, err := tx.Exec(query, applicant.UserID)
	return err
}

func (r *userRepository) createRecruiterTx(tx *sql.Tx, recruiter models.Recruiter) error {
	query := `
		INSERT INTO recruiters (user_id, company_id)
		VALUES ($1, $2)
	`
	_, err := tx.Exec(query, recruiter.UserID, recruiter.CompanyID)
	return err
}

