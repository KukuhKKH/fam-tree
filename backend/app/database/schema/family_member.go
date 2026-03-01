package schema

import "time"

type FamilyMember struct {
	ID       uint64     `gorm:"primaryKey;column:id" json:"id"`
	FamilyID uint64     `gorm:"column:family_id;not null" json:"family_id"`
	UserID   uint64     `gorm:"column:user_id;not null" json:"user_id"`
	Role     string     `gorm:"column:role;not null" json:"role"`
	JoinedAt *time.Time `gorm:"column:joined_at;autoCreateTime" json:"joined_at"`

	// Relationships
	Family Family `gorm:"foreignKey:FamilyID" json:"family,omitempty"`
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
