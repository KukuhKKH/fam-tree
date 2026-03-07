package repository

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-silsilah/internal/bootstrap/database"
)

type RelationshipRepository interface {
	Create(rel *schema.Relationship) (*schema.Relationship, error)
	FindByID(id, familyID uint64) (*schema.Relationship, error)
	ListByFamilyID(familyID uint64) ([]schema.Relationship, error)
	ListByPersonID(personID, familyID uint64) ([]schema.Relationship, error)
	Delete(id, familyID uint64) error
	ExistsDuplicate(familyID, personAID, personBID uint64, relType string) bool
	Update(rel *schema.Relationship) error
}

type relationshipRepository struct {
	DB *database.Database
}

func NewRelationshipRepository(db *database.Database) RelationshipRepository {
	return &relationshipRepository{DB: db}
}

func (r *relationshipRepository) Create(rel *schema.Relationship) (*schema.Relationship, error) {
	if err := r.DB.DB.Create(rel).Error; err != nil {
		return nil, err
	}
	return rel, nil
}

func (r *relationshipRepository) FindByID(id, familyID uint64) (*schema.Relationship, error) {
	var rel schema.Relationship
	if err := r.DB.DB.
		Where("id = ? AND family_id = ?", id, familyID).
		First(&rel).Error; err != nil {
		return nil, err
	}
	return &rel, nil
}

func (r *relationshipRepository) ListByFamilyID(familyID uint64) ([]schema.Relationship, error) {
	var rels []schema.Relationship
	if err := r.DB.DB.
		Where("family_id = ?", familyID).
		Order("created_at ASC").
		Find(&rels).Error; err != nil {
		return nil, err
	}
	return rels, nil
}

// ListByPersonID returns all relationships where given person appears on either side.
func (r *relationshipRepository) ListByPersonID(personID, familyID uint64) ([]schema.Relationship, error) {
	var rels []schema.Relationship
	if err := r.DB.DB.
		Where("family_id = ? AND (person_a_id = ? OR person_b_id = ?)", familyID, personID, personID).
		Find(&rels).Error; err != nil {
		return nil, err
	}
	return rels, nil
}

func (r *relationshipRepository) Delete(id, familyID uint64) error {
	return r.DB.DB.
		Where("id = ? AND family_id = ?", id, familyID).
		Delete(&schema.Relationship{}).Error
}

// ExistsDuplicate checks for symmetric duplicate relationships.
// For directed (parent_child): checks same direction only.
// For symmetric (spouse, sibling): checks both (A,B) and (B,A).
func (r *relationshipRepository) ExistsDuplicate(familyID, personAID, personBID uint64, relType string) bool {
	var count int64

	switch relType {
	case schema.RelationshipParentChild:
		// Directed: only same direction matters
		r.DB.DB.Model(&schema.Relationship{}).
			Where("family_id = ? AND person_a_id = ? AND person_b_id = ? AND relationship_type = ?",
				familyID, personAID, personBID, relType).
			Count(&count)
	default:
		// Symmetric: (A,B) same as (B,A)
		r.DB.DB.Model(&schema.Relationship{}).
			Where(`family_id = ? AND relationship_type = ? AND (
				(person_a_id = ? AND person_b_id = ?) OR
				(person_a_id = ? AND person_b_id = ?)
			)`, familyID, relType, personAID, personBID, personBID, personAID).
			Count(&count)
	}

	return count > 0
}

func (r *relationshipRepository) Update(rel *schema.Relationship) error {
	return r.DB.DB.Save(rel).Error
}
