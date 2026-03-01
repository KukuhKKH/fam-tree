package schema

import (
	"time"

	"git.dev.siap.id/kukuhkkh/app-silsilah/utils/helpers"
)

type User struct {
	ID        uint64     `gorm:"primary_key;column:id" json:"id"`
	Name      string     `gorm:"column:name;not null" json:"name"`
	Password  *string    `gorm:"column:password" json:"-"`
	Email     string     `gorm:"column:email;unique;not null" json:"email"`
	LogtoSub  string     `gorm:"type:varchar(255);column:logto_sub;uniqueIndex;not null" json:"logto_sub"`
	LastLogin *time.Time `gorm:"column:last_login" json:"last_login"`
	Roles     []Role     `gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID" json:"roles,omitempty"`
	Base
}

// ComparePassword compare password
func (u *User) ComparePassword(password string) bool {
	if u.Password == nil {
		return false
	}

	return helpers.ValidateHash(password, *u.Password)
}
