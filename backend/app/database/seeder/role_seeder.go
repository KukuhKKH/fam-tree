package seeder

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	"gorm.io/gorm"
)

// RoleSeeder seeds the default roles.
type RoleSeeder struct{}

func (s *RoleSeeder) Count() (int, error) {
	return 0, nil
}

// Seed inserts default roles if they don't exist.
func (s *RoleSeeder) Seed(db *gorm.DB) error {
	defaults := []schema.Role{
		{Name: "Super Admin", Code: "super_admin", Description: "Full system access"},
		{Name: "Admin Account", Code: "admin_account", Description: "Full system access"},
		{Name: "User", Code: "user", Description: "Regular user"},
	}

	for _, role := range defaults {
		if err := db.Where(schema.Role{Code: role.Code}).
			FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}

	return nil
}
