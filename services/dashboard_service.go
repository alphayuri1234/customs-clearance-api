package services

import (
	"customs-clearance-api/models"
	"gorm.io/gorm"
)

type DashboardService struct {
	db *gorm.DB
}

func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

func (s *DashboardService) GetSummary() (models.DashboardSummary, error) {
	var summary models.DashboardSummary

	// 1. Total Clearances
	if err := s.db.Model(&models.Clearance{}).Count(&summary.TotalClearances).Error; err != nil {
		return summary, err
	}

	// 2. Total Valuation
	if err := s.db.Model(&models.Clearance{}).Select("COALESCE(SUM(valuation), 0)").Scan(&summary.TotalValuation).Error; err != nil {
		return summary, err
	}

	// 3. Status Counts
	var statusResults []struct {
		Status string
		Count  int64
	}
	if err := s.db.Model(&models.Clearance{}).Select("status, COUNT(*) as count").Group("status").Scan(&statusResults).Error; err != nil {
		return summary, err
	}
	summary.StatusCounts = make(map[string]int64)
	// Inisialisasi semua kemungkinan status agar terdefinisi di response
	allStatuses := []string{
		models.StatusSubmitted, models.StatusInspection, models.StatusInspectionPassed,
		models.StatusApproved, models.StatusReleased, models.StatusHold, models.StatusGateOut,
	}
	for _, status := range allStatuses {
		summary.StatusCounts[status] = 0
	}
	for _, r := range statusResults {
		summary.StatusCounts[r.Status] = r.Count
	}

	// 4. Risk Level Counts
	var riskResults []struct {
		Level string
		Count int64
	}
	if err := s.db.Model(&models.RiskProfile{}).Select("level, COUNT(*) as count").Group("level").Scan(&riskResults).Error; err != nil {
		return summary, err
	}
	summary.RiskCounts = make(map[string]int64)
	summary.RiskCounts[models.RiskLevelHigh] = 0
	summary.RiskCounts[models.RiskLevelLow] = 0
	for _, r := range riskResults {
		summary.RiskCounts[r.Level] = r.Count
	}

	// 5. Top 3 Commodities
	if err := s.db.Table("clearances").
		Select("commodities.hs_code, commodities.description, COUNT(*) as count, SUM(clearances.valuation) as total_value").
		Joins("JOIN commodities ON commodities.id = clearances.commodity_id").
		Group("commodities.hs_code, commodities.description").
		Order("count DESC, total_value DESC").
		Limit(3).
		Scan(&summary.TopCommodities).Error; err != nil {
		return summary, err
	}

	// 6. Top 3 Ports
	if err := s.db.Table("clearances").
		Select("ports.code, ports.name, COUNT(*) as count, SUM(clearances.valuation) as total_value").
		Joins("JOIN ports ON ports.id = clearances.port_id").
		Group("ports.code, ports.name").
		Order("count DESC, total_value DESC").
		Limit(3).
		Scan(&summary.TopPorts).Error; err != nil {
		return summary, err
	}

	// 7. Recent 5 Clearances
	if err := s.db.Preload("User").
		Preload("Commodity").
		Preload("Port").
		Preload("Port.Country").
		Preload("RiskProfile").
		Order("id DESC").
		Limit(5).
		Find(&summary.RecentClearances).Error; err != nil {
		return summary, err
	}

	return summary, nil
}
