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
		UserID:   user.ID,
		FullName: *req.FullName,
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

func (r *userRepository) createApplicantTx(tx *sql.Tx, applicant models.Applicant) error {
	query := `
		INSERT INTO applicants (user_id, full_name)
		VALUES ($1, $2)
	`
	_, err := tx.Exec(query, applicant.UserID, applicant.FullName)
	return err
}
