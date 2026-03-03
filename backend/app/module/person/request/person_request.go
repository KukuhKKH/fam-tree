package request

import "encoding/json"

type CreatePersonRequest struct {
	FullName   string  `json:"full_name"   validate:"required,min=2"`
	Nickname   string  `json:"nickname"`
	Gender     string  `json:"gender"      validate:"required,oneof=male female other"`
	BirthDate  string  `json:"birth_date"` // RFC3339 / YYYY-MM-DD, optional
	BirthPlace string  `json:"birth_place"`
	IsAlive    *bool   `json:"is_alive"    validate:"required"`
	DeathDate  string  `json:"death_date"` // required if is_alive=false
	DeathPlace string  `json:"death_place"`
	PhotoURL   string  `json:"photo_url"`
	Bio        string  `json:"bio"`
	UserID     *uint64 `json:"user_id"` // optional: link to registered user
}

type UpdatePersonRequest struct {
	FullName   string  `json:"full_name"   validate:"omitempty,min=2"`
	Nickname   string  `json:"nickname"`
	Gender     string  `json:"gender"      validate:"omitempty,oneof=male female other"`
	BirthDate  string  `json:"birth_date"`
	BirthPlace string  `json:"birth_place"`
	IsAlive    *bool   `json:"is_alive"`
	DeathDate  string  `json:"death_date"`
	DeathPlace string  `json:"death_place"`
	PhotoURL   string  `json:"photo_url"`
	Bio        string  `json:"bio"`
	UserID     *uint64 `json:"user_id"`
}

type CreateRelationshipRequest struct {
	PersonAID        uint64          `json:"person_a_id"        validate:"required"`
	PersonBID        uint64          `json:"person_b_id"        validate:"required"`
	RelationshipType string          `json:"relationship_type"  validate:"required,oneof=parent_child spouse sibling"`
	Metadata         json.RawMessage `json:"metadata"` // optional, e.g. {"marriage_date":"2020-01-01"}
}
