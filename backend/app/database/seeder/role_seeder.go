package seeder

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	"gorm.io/gorm"
)

type RoleSeeder struct {
	db *gorm.DB
}

func NewRoleSeeder(db *gorm.DB) *RoleSeeder {
	return &RoleSeeder{db: db}
}

func (s *RoleSeeder) Count() (int, error) {
	if s.db == nil {
		return 0, nil
	}

	var count int64
	err := s.db.Model(&schema.Role{}).Count(&count).Error
	return int(count), err
}

// Seed does nothing - roles come from Logto, not from seeding.
func (s *RoleSeeder) Seed(_ *gorm.DB) error {
	return nil
}
