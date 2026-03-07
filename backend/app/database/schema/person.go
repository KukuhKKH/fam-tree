package schema

import "database/sql"

const (
	GenderMale   = "male"
	GenderFemale = "female"
	GenderOther  = "other"
)

type Person struct {
	ID         uint64       `gorm:"primaryKey;column:id"                                       json:"id"`
	FamilyID   uint64       `gorm:"column:family_id;not null;index"                            json:"family_id"`
	FullName   string       `gorm:"column:full_name;not null"                                  json:"full_name"`
	Nickname   string       `gorm:"column:nickname"                                            json:"nickname"`
	Gender     string       `gorm:"column:gender;not null"                                     json:"gender"` // male, female, other
	BirthDate  sql.NullTime `gorm:"column:birth_date"                                          json:"birth_date"`
	BirthPlace string       `gorm:"column:birth_place"                                         json:"birth_place"`
	IsAlive    *bool        `gorm:"column:is_alive;not null;default:true"                      json:"is_alive"`
	DeathDate  sql.NullTime `gorm:"column:death_date"                                          json:"death_date"`
	DeathPlace string       `gorm:"column:death_place"                                         json:"death_place"`
	PhotoURL   string       `gorm:"column:photo_url"                                           json:"photo_url"`
	Bio        string       `gorm:"column:bio;type:text"                                       json:"bio"`
	UserID     *uint64      `gorm:"column:user_id;index"                                       json:"user_id"` // nullable, linked account

	Family Family `gorm:"foreignKey:FamilyID"                                        json:"family,omitempty"`
	User   *User  `gorm:"foreignKey:UserID"                                          json:"user,omitempty"`
	Base
}
