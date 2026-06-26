package models

import "time"

type InspectionResult struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ClearanceID uint      `json:"clearance_id" gorm:"uniqueIndex;not null"`
	OfficerID   uint      `json:"officer_id" gorm:"not null"`
	Officer     Officer   `json:"officer" gorm:"foreignKey:OfficerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Result      string    `json:"result" gorm:"not null"` // PASS / FAIL
	Notes       string    `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
