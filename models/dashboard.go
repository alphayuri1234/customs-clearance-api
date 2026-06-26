package models

type TopCommodity struct {
	HSCode      string  `json:"hs_code"`
	Description string  `json:"description"`
	Count       int64   `json:"count"`
	TotalValue  float64 `json:"total_value"`
}

type TopPort struct {
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Count       int64   `json:"count"`
	TotalValue  float64 `json:"total_value"`
}

type DashboardSummary struct {
	TotalClearances  int64             `json:"total_clearances"`
	StatusCounts     map[string]int64  `json:"status_counts"`
	RiskCounts       map[string]int64  `json:"risk_counts"`
	TotalValuation   float64           `json:"total_valuation"`
	TopCommodities   []TopCommodity    `json:"top_commodities"`
	TopPorts         []TopPort         `json:"top_ports"`
	RecentClearances []Clearance       `json:"recent_clearances"`
}
