package repository

import (
	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-silsilah/internal/bootstrap/database"
)

type PersonRepository interface {
	Create(person *schema.Person) (*schema.Person, error)
	FindByID(id, familyID uint64) (*schema.Person, error)
	ListByFamilyID(familyID uint64) ([]schema.Person, error)
	Update(person *schema.Person) error
	Delete(id, familyID uint64) error

	// Tree traversal — uses Recursive CTE for optimal performance
	GetAncestors(personID, familyID uint64) ([]schema.Person, error)
	GetDescendants(personID, familyID uint64) ([]schema.Person, error)
}

type personRepository struct {
	DB *database.Database
}

func NewPersonRepository(db *database.Database) PersonRepository {
	return &personRepository{DB: db}
}

func (r *personRepository) Create(person *schema.Person) (*schema.Person, error) {
	if err := r.DB.DB.Create(person).Error; err != nil {
		return nil, err
	}
	return person, nil
}

func (r *personRepository) FindByID(id, familyID uint64) (*schema.Person, error) {
	var person schema.Person
	if err := r.DB.DB.
		Where("id = ? AND family_id = ?", id, familyID).
		First(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *personRepository) ListByFamilyID(familyID uint64) ([]schema.Person, error) {
	var persons []schema.Person
	if err := r.DB.DB.
		Where("family_id = ?", familyID).
		Order("full_name ASC").
		Find(&persons).Error; err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *personRepository) Update(person *schema.Person) error {
	return r.DB.DB.Save(person).Error
}

func (r *personRepository) Delete(id, familyID uint64) error {
	return r.DB.DB.
		Where("id = ? AND family_id = ?", id, familyID).
		Delete(&schema.Person{}).Error
}

func (r *personRepository) GetAncestors(personID, familyID uint64) ([]schema.Person, error) {
	const query = `
		WITH RECURSIVE ancestors AS (
			-- Base: direct parents of the given person
			SELECT r.person_a_id AS id
			FROM relationships r
			WHERE r.person_b_id = ?
			  AND r.family_id   = ?
			  AND r.relationship_type = 'parent_child'
			  AND r.deleted_at IS NULL

			UNION ALL

			-- Recursive: parents of the parents
			SELECT r.person_a_id
			FROM relationships r
			INNER JOIN ancestors a ON r.person_b_id = a.id
			WHERE r.family_id = ?
			  AND r.relationship_type = 'parent_child'
			  AND r.deleted_at IS NULL
		)
		SELECT p.* FROM persons p
		INNER JOIN ancestors a ON p.id = a.id
		WHERE p.deleted_at IS NULL
		ORDER BY p.full_name ASC
	`
	var persons []schema.Person
	if err := r.DB.DB.Raw(query, personID, familyID, familyID).Scan(&persons).Error; err != nil {
		return nil, err
	}
	return persons, nil
}

func (r *personRepository) GetDescendants(personID, familyID uint64) ([]schema.Person, error) {
	const query = `
		WITH RECURSIVE descendants AS (
			-- Base: direct children of the given person
			SELECT r.person_b_id AS id
			FROM relationships r
			WHERE r.person_a_id = ?
			  AND r.family_id   = ?
			  AND r.relationship_type = 'parent_child'
			  AND r.deleted_at IS NULL

			UNION ALL

			-- Recursive: children of the children
			SELECT r.person_b_id
			FROM relationships r
			INNER JOIN descendants d ON r.person_a_id = d.id
			WHERE r.family_id = ?
			  AND r.relationship_type = 'parent_child'
			  AND r.deleted_at IS NULL
		)
		SELECT p.* FROM persons p
		INNER JOIN descendants d ON p.id = d.id
		WHERE p.deleted_at IS NULL
		ORDER BY p.full_name ASC
	`
	var persons []schema.Person
	if err := r.DB.DB.Raw(query, personID, familyID, familyID).Scan(&persons).Error; err != nil {
		return nil, err
	}
	return persons, nil
}
