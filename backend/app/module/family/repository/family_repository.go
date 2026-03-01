package repository

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-silsilah/internal/bootstrap/database"
)

type FamilyRepository interface {
	Create(family *schema.Family) (*schema.Family, error)
	FindBySlug(slug string) (*schema.Family, error)
	FindByID(id uint64) (*schema.Family, error)
	Update(family *schema.Family) error
	Delete(id uint64) error
	GetFamiliesByUserID(userID uint64) ([]schema.FamilyMember, error)
	CheckSlugExists(slug string) bool
}

type familyRepository struct {
	DB *database.Database
}

func NewFamilyRepository(db *database.Database) FamilyRepository {
	return &familyRepository{DB: db}
}

func (r *familyRepository) Create(family *schema.Family) (*schema.Family, error) {
	if err := r.DB.DB.Create(family).Error; err != nil {
		return nil, err
	}
	return family, nil
}

func (r *familyRepository) FindBySlug(slug string) (*schema.Family, error) {
	var family schema.Family
	if err := r.DB.DB.Where("slug = ?", slug).First(&family).Error; err != nil {
		return nil, err
	}
	return &family, nil
}

func (r *familyRepository) FindByID(id uint64) (*schema.Family, error) {
	var family schema.Family
	if err := r.DB.DB.First(&family, id).Error; err != nil {
		return nil, err
	}

	return &family, nil
}

func (r *familyRepository) Update(family *schema.Family) error {
	return r.DB.DB.Save(family).Error
}

func (r *familyRepository) Delete(id uint64) error {
	return r.DB.DB.Delete(&schema.Family{}, id).Error
}

func (r *familyRepository) GetFamiliesByUserID(userID uint64) ([]schema.FamilyMember, error) {
	var members []schema.FamilyMember
	if err := r.DB.DB.Preload("Family").Where("user_id = ?", userID).Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (r *familyRepository) CheckSlugExists(slug string) bool {
	var count int64
	r.DB.DB.Model(&schema.Family{}).Where("slug = ?", slug).Count(&count)

	return count > 0
}
