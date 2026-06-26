package models

import "time"

type ReleaseOrder struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ClearanceID uint      `json:"clearance_id" gorm:"uniqueIndex;not null"`
	ReleaseNo   string    `json:"release_no" gorm:"uniqueIndex;not null"`
	OfficerID   uint      `json:"officer_id" gorm:"not null"`
	Officer     Officer   `json:"officer" gorm:"foreignKey:OfficerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	IssuedAt    time.Time `json:"issued_at" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
