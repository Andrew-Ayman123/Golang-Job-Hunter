// Updated repository/user.go
package repository

import (
	"database/sql"
	"fmt"

	"github.com/Andrew-Ayman123/Job-Hunter/internal/dto"
	"github.com/Andrew-Ayman123/Job-Hunter/internal/models"
)

type AdminRepository interface {
	CreateAdmin(req dto.CreateAdminRequest) (*models.Admin, error)
	CreateRecruiter(req dto.CreateRecruiterRequest) (*models.Recruiter, error)
	CreateCompany(req dto.CreateCompanyRequest) (*models.Company, error)
	UpdateCompany(id string, req dto.UpdateCompanyRequest) (*models.Company, error)
	DeleteCompany(id string) error
}

type adminRepository struct {
	db       *sql.DB
	userRepo UserRepository // Add UserRepository as a dependency
}

func NewAdminRepository(db *sql.DB, userRepo UserRepository) AdminRepository {
	return &adminRepository{
		db:       db,
		userRepo: userRepo,
	}
}

func (r *adminRepository) CreateAdmin(req dto.CreateAdminRequest) (*models.Admin, error) {

	user, err := r.userRepo.CreateUser(req.CreateUserRequest, "admin")
	if err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}
	 
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create admin profile
	admin := models.Admin{
		User: *user,
		AdminLevel: req.AdminLevel,
	}
	err = r.createAdminTx(tx, admin)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &admin, nil
}

func (r *adminRepository) CreateRecruiter(req dto.CreateRecruiterRequest) (*models.Recruiter, error) {

	user, err := r.userRepo.CreateUser(req.CreateUserRequest, "recruiter")
	if err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}
	 
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Create recruiter profile
	recruiter := models.Recruiter{
		User:     *user,
		CompanyID: req.CompanyID,
	}
	err = r.createRecruiterTx(tx, recruiter)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &recruiter, nil
}

func (r *adminRepository) CreateCompany(req dto.CreateCompanyRequest) (*models.Company, error) {
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

func (r *adminRepository) UpdateCompany(id string, req dto.UpdateCompanyRequest) (*models.Company, error) {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Validate company ID
	if id == "" {
		return nil, fmt.Errorf("company ID is required")
	}

	// Check if company exists
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM companies WHERE id = $1)`
	err = tx.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if company exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("company with ID %s not found", id)
	}

	// Prepare update query
	query = `
		UPDATE companies
		SET name = COALESCE($1, name), description = COALESCE($2, description), updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`
	var company models.Company
	err = tx.QueryRow(query, req.Name, req.Description, id).Scan(
		&company.ID, &company.Name, &company.Description, &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("company with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &company, nil
}

func (r *adminRepository) DeleteCompany(id string) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Validate company ID
	if id == "" {
		return fmt.Errorf("company ID is required")
	}

	// Check if company exists
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM companies WHERE id = $1)`
	err = tx.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if company exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("company with ID %s not found", id)
	}

	// Delete company
	query = `DELETE FROM companies WHERE id = $1`
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *adminRepository) createRecruiterTx(tx *sql.Tx, recruiter models.Recruiter) error {
	query := `
		INSERT INTO recruiters (user_id, company_id)
		VALUES ($1, $2)
	`
	_, err := tx.Exec(query, recruiter.ID, recruiter.CompanyID)
	return err
}

func (r *adminRepository) createAdminTx(tx *sql.Tx, admin models.Admin) error {
	query := `
		INSERT INTO admins (user_id, admin_level)
		VALUES ($1, $2)
	`
	_, err := tx.Exec(query, admin.ID, admin.AdminLevel)
	return err
}
