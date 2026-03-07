package schema

import (
	"database/sql"
	"encoding/json"
)

const (
	RelationshipParentChild = "parent_child"
	RelationshipSpouse      = "spouse"
	RelationshipSibling     = "sibling"
)

type Relationship struct {
	ID               uint64          `gorm:"primaryKey;column:id"                                       json:"id"`
	FamilyID         uint64          `gorm:"column:family_id;not null;index:idx_rel_family"             json:"family_id"`
	PersonAID        uint64          `gorm:"column:person_a_id;not null;index:idx_rel_family"           json:"person_a_id"`
	PersonBID        uint64          `gorm:"column:person_b_id;not null;index:idx_rel_family"           json:"person_b_id"`
	RelationshipType string          `gorm:"column:relationship_type;not null"                          json:"relationship_type"`
	MarriageDate     sql.NullTime    `gorm:"column:marriage_date"                                       json:"marriage_date"`
	DivorceDate      sql.NullTime    `gorm:"column:divorce_date"                                        json:"divorce_date"`
	Metadata         json.RawMessage `gorm:"column:metadata;type:jsonb;default:'{}'"                    json:"metadata"`

	Family  Family `gorm:"foreignKey:FamilyID"  json:"family,omitempty"`
	PersonA Person `gorm:"foreignKey:PersonAID" json:"person_a,omitempty"`
	PersonB Person `gorm:"foreignKey:PersonBID" json:"person_b,omitempty"`
	Base
}
