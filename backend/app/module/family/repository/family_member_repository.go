package repository

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-silsilah/internal/bootstrap/database"
)

type FamilyMemberRepository interface {
	Create(member *schema.FamilyMember) (*schema.FamilyMember, error)
	FindByFamilyAndUser(familyID, userID uint64) (*schema.FamilyMember, error)
	GetByFamilyID(familyID uint64) ([]schema.FamilyMember, error)
	UpdateRole(familyID, userID uint64, newRole string) error
	Delete(familyID, userID uint64) error
	CheckUserRoleInFamily(familyID, userID uint64) (string, error)
}

type familyMemberRepository struct {
	DB *database.Database
}

func NewFamilyMemberRepository(db *database.Database) FamilyMemberRepository {
	return &familyMemberRepository{DB: db}
}

func (r *familyMemberRepository) Create(member *schema.FamilyMember) (*schema.FamilyMember, error) {
	if err := r.DB.DB.Create(member).Error; err != nil {
		return nil, err
	}
	return member, nil
}

func (r *familyMemberRepository) FindByFamilyAndUser(familyID, userID uint64) (*schema.FamilyMember, error) {
	var member schema.FamilyMember
	if err := r.DB.DB.Where("family_id = ? AND user_id = ?", familyID, userID).First(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

func (r *familyMemberRepository) GetByFamilyID(familyID uint64) ([]schema.FamilyMember, error) {
	var members []schema.FamilyMember
	if err := r.DB.DB.Preload("User").Where("family_id = ?", familyID).Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

func (r *familyMemberRepository) UpdateRole(familyID, userID uint64, newRole string) error {
	return r.DB.DB.Model(&schema.FamilyMember{}).
		Where("family_id = ? AND user_id = ?", familyID, userID).
		Update("role", newRole).Error
}

func (r *familyMemberRepository) Delete(familyID, userID uint64) error {
	return r.DB.DB.Where("family_id = ? AND user_id = ?", familyID, userID).Delete(&schema.FamilyMember{}).Error
}

func (r *familyMemberRepository) CheckUserRoleInFamily(familyID, userID uint64) (string, error) {
	var role string
	err := r.DB.DB.Model(&schema.FamilyMember{}).
		Where("family_id = ? AND user_id = ?", familyID, userID).
		Select("role").Row().Scan(&role)

	return role, err
}
