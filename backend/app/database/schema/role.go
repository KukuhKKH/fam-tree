package schema

type Role struct {
	ID          uint64 `gorm:"primaryKey;column:id" json:"id"`
	Name        string `gorm:"column:name;not null;uniqueIndex" json:"name"`
	Code        string `gorm:"column:code;not null;uniqueIndex" json:"code"`
	Description string `gorm:"column:description" json:"description"`
	Base
}
