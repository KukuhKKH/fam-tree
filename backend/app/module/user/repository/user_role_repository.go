package repository

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-silsilah/internal/bootstrap/database"
	"gorm.io/gorm"
)

type UserRoleRepository interface {
	GetUserRoles(userID uint64) ([]schema.Role, error)
	SyncUserRoles(userID uint64, roles []schema.Role) error
}

type userRoleRepository struct {
	DB *database.Database
}

func NewUserRoleRepository(db *database.Database) UserRoleRepository {
	return &userRoleRepository{DB: db}
}

func (r *userRoleRepository) GetUserRoles(userID uint64) ([]schema.Role, error) {
	var roles []schema.Role
	err := r.DB.DB.
		Joins("JOIN user_roles ur ON ur.role_id = roles.id").
		Where("ur.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// SyncUserRoles replaces the user's roles atomically.
// It deletes existing user_roles rows and re-inserts the new set.
func (r *userRoleRepository) SyncUserRoles(userID uint64, roles []schema.Role) error {
	return r.DB.DB.Transaction(func(tx *gorm.DB) error {
		// delete existing
		if err := tx.Where("user_id = ?", userID).Delete(&schema.UserRole{}).Error; err != nil {
			return err
		}

		// insert new
		for _, role := range roles {
			ur := schema.UserRole{UserID: userID, RoleID: role.ID}
			if err := tx.Create(&ur).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
