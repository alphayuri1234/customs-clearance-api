package models

import "time"

const (
	RiskLevelHigh = "HIGH"
	RiskLevelLow  = "LOW"
)

type RiskProfile struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ClearanceID uint      `json:"clearance_id" gorm:"uniqueIndex;not null"`
	Level       string    `json:"level" gorm:"not null;default:LOW"` // HIGH or LOW
	Score       float64   `json:"score" gorm:"not null;default:0"`
	Reason      string    `json:"reason" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
