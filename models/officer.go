package models

type Officer struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   uint   `json:"user_id" gorm:"uniqueIndex;not null"`
	User     User   `json:"user" gorm:"foreignKey:UserID"`
	NIP      string `json:"nip" gorm:"uniqueIndex"`
	Position string `json:"position"`
}
