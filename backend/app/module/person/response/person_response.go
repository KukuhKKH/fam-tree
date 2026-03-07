package response

import (
	"encoding/json"
	"time"

	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
)

type PersonResponse struct {
	ID         uint64     `json:"id"`
	FamilyID   uint64     `json:"family_id"`
	FullName   string     `json:"full_name"`
	Nickname   string     `json:"nickname"`
	Gender     string     `json:"gender"`
	BirthDate  *time.Time `json:"birth_date"`
	BirthPlace string     `json:"birth_place"`
	IsAlive    bool       `json:"is_alive"`
	DeathDate  *time.Time `json:"death_date"`
	DeathPlace string     `json:"death_place"`
	PhotoURL   string     `json:"photo_url"`
	Bio        string     `json:"bio"`
	UserID     *uint64    `json:"user_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type RelationshipResponse struct {
	ID               uint64          `json:"id"`
	FamilyID         uint64          `json:"family_id"`
	PersonAID        uint64          `json:"person_a_id"`
	PersonBID        uint64          `json:"person_b_id"`
	RelationshipType string          `json:"relationship_type"`
	MarriageDate     *time.Time      `json:"marriage_date,omitempty"`
	DivorceDate      *time.Time      `json:"divorce_date,omitempty"`
	Metadata         json.RawMessage `json:"metadata"`
	CreatedAt        time.Time       `json:"created_at"`
}

// TreeNode is a single person node for frontend visualization (Cytoscape/D3 compatible)
type TreeNode struct {
	ID         uint64     `json:"id"`
	FullName   string     `json:"full_name"`
	Nickname   string     `json:"nickname"`
	Gender     string     `json:"gender"`
	BirthDate  *time.Time `json:"birth_date"`
	BirthPlace string     `json:"birth_place"`
	IsAlive    bool       `json:"is_alive"`
	DeathDate  *time.Time `json:"death_date"`
	DeathPlace string     `json:"death_place"`
	PhotoURL   string     `json:"photo_url"`
	Bio        string     `json:"bio"`
	UserID     *uint64    `json:"user_id"`
}

// TreeEdge is a relationship edge for frontend visualization
type TreeEdge struct {
	ID               uint64          `json:"id"`
	From             uint64          `json:"from"`
	To               uint64          `json:"to"`
	RelationshipType string          `json:"relationship_type"`
	MarriageDate     *time.Time      `json:"marriage_date,omitempty"`
	DivorceDate      *time.Time      `json:"divorce_date,omitempty"`
	Metadata         json.RawMessage `json:"metadata"`
}

// TreeResponse is the full tree payload sent to the frontend
type TreeResponse struct {
	Nodes  []TreeNode `json:"nodes"`
	Edges  []TreeEdge `json:"edges"`
	RootID uint64     `json:"root_id"`
}

func FromPersonSchema(p *schema.Person) PersonResponse {
	resp := PersonResponse{
		ID:         p.ID,
		FamilyID:   p.FamilyID,
		FullName:   p.FullName,
		Nickname:   p.Nickname,
		Gender:     p.Gender,
		BirthPlace: p.BirthPlace,
		DeathPlace: p.DeathPlace,
		PhotoURL:   p.PhotoURL,
		Bio:        p.Bio,
		UserID:     p.UserID,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}

	if p.IsAlive != nil {
		resp.IsAlive = *p.IsAlive
	} else {
		resp.IsAlive = true
	}

	if p.BirthDate.Valid {
		resp.BirthDate = &p.BirthDate.Time
	}
	if p.DeathDate.Valid {
		resp.DeathDate = &p.DeathDate.Time
	}
	return resp
}

func FromRelationshipSchema(r *schema.Relationship) RelationshipResponse {
	resp := RelationshipResponse{
		ID:               r.ID,
		FamilyID:         r.FamilyID,
		PersonAID:        r.PersonAID,
		PersonBID:        r.PersonBID,
		RelationshipType: r.RelationshipType,
		Metadata:         r.Metadata,
		CreatedAt:        r.CreatedAt,
	}

	if r.MarriageDate.Valid {
		resp.MarriageDate = &r.MarriageDate.Time
	}
	if r.DivorceDate.Valid {
		resp.DivorceDate = &r.DivorceDate.Time
	}

	return resp
}
