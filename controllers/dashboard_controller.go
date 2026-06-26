package controllers

import (
	"net/http"

	"customs-clearance-api/models"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	seederService    *services.SeederService
	dashboardService *services.DashboardService
}

func NewDashboardController(seederService *services.SeederService, dashboardService *services.DashboardService) *DashboardController {
	return &DashboardController{
		seederService:    seederService,
		dashboardService: dashboardService,
	}
}

// SeedData memicu seeding database dengan data dummy Bea Cukai yang representatif
func (c *DashboardController) SeedData(ctx *gin.Context) {
	if err := c.seederService.Seed(); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("gagal melakukan seeding data", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("database berhasil di-seed dengan data sampel Bea Cukai", nil))
}

// GetSummary mengambil agregasi statistik untuk dashboard Bea Cukai
func (c *DashboardController) GetSummary(ctx *gin.Context) {
	summary, err := c.dashboardService.GetSummary()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("gagal mengambil data summary dashboard", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("summary dashboard berhasil diambil", summary))
}
