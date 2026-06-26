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

// SeedData godoc
// @Summary Seed Data Tiruan
// @Description Membersihkan seluruh database (TRUNCATE CASCADE) dan mengisi data tiruan Bea Cukai dalam jumlah besar untuk keperluan pengetesan
// @Tags Dashboard / Seeder
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse "Seeding Berhasil"
// @Failure 500 {object} models.APIResponse "Gagal Melakukan Seeding"
// @Router /seed [post]
func (c *DashboardController) SeedData(ctx *gin.Context) {
	if err := c.seederService.Seed(); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("gagal melakukan seeding data", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("database berhasil di-seed dengan data sampel Bea Cukai", nil))
}

// GetSummary godoc
// @Summary Ringkasan Dashboard Bea Cukai
// @Description Menampilkan statistik pengajuan barang, kategori risiko, top komoditas, top pelabuhan, dan transaksi terbaru (Khusus Officer)
// @Tags Dashboard / Seeder
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.DashboardSummary} "Dashboard Berhasil Diambil"
// @Failure 500 {object} models.APIResponse "Gagal Mengambil Data"
// @Router /dashboard [get]
func (c *DashboardController) GetSummary(ctx *gin.Context) {
	summary, err := c.dashboardService.GetSummary()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("gagal mengambil data summary dashboard", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("summary dashboard berhasil diambil", summary))
}
