package response

import (
	"time"

	"git.dev.siap.id/kukuhkkh/app-silsilah/app/database/schema"
)

type FamilyResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Visibility  string    `json:"visibility"`
	CreatedBy   uint64    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	MyRole      string    `json:"my_role,omitempty"`
}

type FamilyMemberResponse struct {
	ID        uint64     `json:"id"`
	FamilyID  uint64     `json:"family_id"`
	UserID    uint64     `json:"user_id"`
	UserName  string     `json:"user_name"`
	UserEmail string     `json:"user_email"`
	Role      string     `json:"role"`
	JoinedAt  *time.Time `json:"joined_at"`
}

func FromFamilySchema(family *schema.Family, role string) FamilyResponse {
	return FamilyResponse{
		ID:          family.ID,
		Name:        family.Name,
		Slug:        family.Slug,
		Description: family.Description,
		Visibility:  family.Visibility,
		CreatedBy:   family.CreatedBy,
		CreatedAt:   family.CreatedAt,
		UpdatedAt:   family.UpdatedAt,
		MyRole:      role,
	}
}

func FromFamilyMemberSchema(member *schema.FamilyMember) FamilyMemberResponse {
	return FamilyMemberResponse{
		ID:        member.ID,
		FamilyID:  member.FamilyID,
		UserID:    member.UserID,
		UserName:  member.User.Name,
		UserEmail: member.User.Email,
		Role:      member.Role,
		JoinedAt:  member.JoinedAt,
	}
}
