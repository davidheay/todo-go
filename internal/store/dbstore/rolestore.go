package dbstore

type Role struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Role  string `gorm:"unique" json:"role"`
}
