package repository

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-silsilah/internal/bootstrap/database"
)

type RoleRepository interface {
	FindByCode(code string) (*schema.Role, error)
	FindByID(id uint64) (*schema.Role, error)
	FindByCodes(codes []string) ([]schema.Role, error)
	CreateIfNotExists(role *schema.Role) (*schema.Role, error)
}

type roleRepository struct {
	DB *database.Database
}

func NewRoleRepository(db *database.Database) RoleRepository {
	return &roleRepository{DB: db}
}

func (r *roleRepository) FindByCode(code string) (*schema.Role, error) {
	var role schema.Role
	if err := r.DB.DB.Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByID(id uint64) (*schema.Role, error) {
	var role schema.Role
	if err := r.DB.DB.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByCodes(codes []string) ([]schema.Role, error) {
	var roles []schema.Role
	if err := r.DB.DB.Where("code IN ?", codes).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// CreateIfNotExists inserts a role only if the code does not exist yet.
func (r *roleRepository) CreateIfNotExists(role *schema.Role) (*schema.Role, error) {
	existing, err := r.FindByCode(role.Code)
	if err == nil {
		return existing, nil
	}
	if err := r.DB.DB.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}
