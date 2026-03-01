package schema

// Constants for Family Member Roles
const (
	FamilyRoleOwner  = "owner"
	FamilyRoleAdmin  = "admin"
	FamilyRoleEditor = "editor"
	FamilyRoleViewer = "viewer"
)

type Family struct {
	ID          uint64 `gorm:"primaryKey;column:id" json:"id"`
	Name        string `gorm:"column:name;not null" json:"name"`
	Slug        string `gorm:"column:slug;uniqueIndex;not null" json:"slug"`
	Description string `gorm:"column:description" json:"description"`
	Visibility  string `gorm:"column:visibility;default:'private';not null" json:"visibility"`
	CreatedBy   uint64 `gorm:"column:created_by;not null" json:"created_by"`

	// Relationships
	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Base
}
