package schema

type UserRole struct {
	UserID uint64 `gorm:"primaryKey;column:user_id" json:"user_id"`
	RoleID uint64 `gorm:"primaryKey;column:role_id" json:"role_id"`

	// Associations
	User User `gorm:"foreignKey:UserID" json:"-"`
	Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}
