package models

import "time"

type Country struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Code      string    `json:"code" gorm:"uniqueIndex;not null;size:3"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Port struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Code      string    `json:"code" gorm:"uniqueIndex;not null;size:10"`
	Name      string    `json:"name" gorm:"not null"`
	CountryID uint      `json:"country_id" gorm:"not null"`
	Country   Country   `json:"country" gorm:"foreignKey:CountryID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Commodity struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	HSCode         string    `json:"hs_code" gorm:"uniqueIndex;not null;size:20"`
	Description    string    `json:"description" gorm:"not null"`
	ImportDutyRate float64   `json:"import_duty_rate" gorm:"not null;default:0"`
	VATRate        float64   `json:"vat_rate" gorm:"not null;default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CountryRequest struct {
	Code string `json:"code" binding:"required,max=3"`
	Name string `json:"name" binding:"required"`
}

type PortRequest struct {
	Code      string `json:"code" binding:"required,max=10"`
	Name      string `json:"name" binding:"required"`
	CountryID uint   `json:"country_id" binding:"required"`
}

type CommodityRequest struct {
	HSCode         string  `json:"hs_code" binding:"required,max=20"`
	Description    string  `json:"description" binding:"required"`
	ImportDutyRate float64 `json:"import_duty_rate" binding:"min=0"`
	VATRate        float64 `json:"vat_rate" binding:"min=0"`
}
