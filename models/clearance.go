package models

import "time"

const (
	StatusSubmitted        = "SUBMITTED"
	StatusInspection       = "INSPECTION"
	StatusInspectionPassed = "INSPECTION_PASSED"
	StatusApproved         = "APPROVED"
	StatusReleased         = "RELEASED"
	StatusHold             = "HOLD"
	StatusGateOut          = "GATE_OUT"
)

type Clearance struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	UserID      uint         `json:"user_id" gorm:"not null"`
	User        User         `json:"user" gorm:"foreignKey:UserID"`
	CommodityID uint         `json:"commodity_id" gorm:"not null"`
	Commodity   Commodity    `json:"commodity" gorm:"foreignKey:CommodityID"`
	PortID      uint         `json:"port_id" gorm:"not null"`
	Port        Port         `json:"port" gorm:"foreignKey:PortID"`
	Valuation   float64      `json:"valuation" gorm:"not null"`
	Description string       `json:"description" gorm:"not null"`
	Status      string            `json:"status" gorm:"not null;default:SUBMITTED"`
	RiskProfile *RiskProfile      `json:"risk_profile,omitempty" gorm:"foreignKey:ClearanceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Inspection  *InspectionResult `json:"inspection,omitempty" gorm:"foreignKey:ClearanceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Release     *ReleaseOrder     `json:"release,omitempty" gorm:"foreignKey:ClearanceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type InspectionRequest struct {
	ClearanceID uint   `json:"clearance_id" binding:"required"`
	Result      string `json:"result" binding:"required,oneof=PASS FAIL"`
}

type ApproveRequest struct {
	ClearanceID uint `json:"clearance_id" binding:"required"`
}

type ReleaseRequest struct {
	ClearanceID uint `json:"clearance_id" binding:"required"`
}

type GateOutRequest struct {
	ClearanceID uint `json:"clearance_id" binding:"required"`
}
