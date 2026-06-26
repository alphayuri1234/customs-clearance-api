package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"customs-clearance-api/models"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MasterController struct {
	masterService *services.MasterService
}

func NewMasterController(masterService *services.MasterService) *MasterController {
	return &MasterController{masterService: masterService}
}

func parseID(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("id tidak valid", nil))
		return 0, false
	}

	return uint(id), true
}

func respondMasterError(ctx *gin.Context, message string, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse(message+" tidak ditemukan", nil))
		return
	}

	ctx.JSON(http.StatusBadRequest, models.ErrorResponse(message+" gagal diproses", err.Error()))
}

// ListCountries godoc
// @Summary Daftar Negara
// @Description Mengambil semua data master negara (Khusus Officer)
// @Tags Master Negara
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.Country} "Berhasil mengambil data"
// @Failure 500 {object} models.APIResponse "Gagal mengambil data"
// @Router /master/countries [get]
func (controller *MasterController) ListCountries(ctx *gin.Context) {
	countries, err := controller.masterService.ListCountries()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("data negara gagal diambil", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("data negara berhasil diambil", countries))
}

// CreateCountry godoc
// @Summary Tambah Negara
// @Description Membuat data master negara baru (Khusus Officer)
// @Tags Master Negara
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CountryRequest true "Payload Negara"
// @Success 201 {object} models.APIResponse{data=models.Country} "Negara Berhasil Dibuat"
// @Failure 400 {object} models.APIResponse "Request Tidak Valid atau Gagal Diproses"
// @Router /master/countries [post]
func (controller *MasterController) CreateCountry(ctx *gin.Context) {
	var request models.CountryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request negara tidak valid", err.Error()))
		return
	}

	country, err := controller.masterService.CreateCountry(request)
	if err != nil {
		respondMasterError(ctx, "negara", err)
		return
	}

	ctx.JSON(http.StatusCreated, models.SuccessResponse("negara berhasil dibuat", country))
}

// GetCountry godoc
// @Summary Detail Negara
// @Description Mengambil satu data master negara berdasarkan ID (Khusus Officer)
// @Tags Master Negara
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Negara ID"
// @Success 200 {object} models.APIResponse{data=models.Country} "Negara Berhasil Diambil"
// @Failure 400 {object} models.APIResponse "ID Tidak Valid"
// @Failure 404 {object} models.APIResponse "Negara Tidak Ditemukan"
// @Router /master/countries/{id} [get]
func (controller *MasterController) GetCountry(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	country, err := controller.masterService.GetCountry(id)
	if err != nil {
		respondMasterError(ctx, "negara", err)
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("negara berhasil diambil", country))
}

// UpdateCountry godoc
// @Summary Perbarui Negara
// @Description Mengubah data master negara berdasarkan ID (Khusus Officer)
// @Tags Master Negara
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Negara ID"
// @Param request body models.CountryRequest true "Payload Negara"
// @Success 200 {object} models.APIResponse{data=models.Country} "Negara Berhasil Diperbarui"
// @Failure 400 {object} models.APIResponse "Request atau ID Tidak Valid"
// @Failure 404 {object} models.APIResponse "Negara Tidak Ditemukan"
// @Router /master/countries/{id} [put]
func (controller *MasterController) UpdateCountry(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var request models.CountryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request negara tidak valid", err.Error()))
		return
	}

	country, err := controller.masterService.UpdateCountry(id, request)
	if err != nil {
		respondMasterError(ctx, "negara", err)
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("negara berhasil diperbarui", country))
}

// DeleteCountry godoc
// @Summary Hapus Negara
// @Description Menghapus data master negara berdasarkan ID (Khusus Officer)
// @Tags Master Negara
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Negara ID"
// @Success 200 {object} models.APIResponse "Negara Berhasil Dihapus"
// @Failure 400 {object} models.APIResponse "ID Tidak Valid"
// @Failure 404 {object} models.APIResponse "Negara Tidak Ditemukan"
// @Router /master/countries/{id} [delete]
func (controller *MasterController) DeleteCountry(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	if err := controller.masterService.DeleteCountry(id); err != nil {
		respondMasterError(ctx, "negara", err)
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("negara berhasil dihapus", nil))
}

// ListPorts godoc
// @Summary Daftar Pelabuhan
// @Description Mengambil semua data master pelabuhan (Khusus Officer)
// @Tags Master Pelabuhan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.Port} "Pelabuhan Berhasil Diambil"
// @Failure 500 {object} models.APIResponse "Gagal Mengambil Data"
// @Router /master/ports [get]
func (controller *MasterController) ListPorts(ctx *gin.Context) {
	ports, err := controller.masterService.ListPorts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("data pelabuhan gagal diambil", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("data pelabuhan berhasil diambil", ports))
}

// CreatePort godoc
// @Summary Tambah Pelabuhan
// @Description Membuat data master pelabuhan baru (Khusus Officer)
// @Tags Master Pelabuhan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.PortRequest true "Payload Pelabuhan"
// @Success 201 {object} models.APIResponse{data=models.Port} "Pelabuhan Berhasil Dibuat"
// @Failure 400 {object} models.APIResponse "Request Tidak Valid atau Gagal Diproses"
// @Router /master/ports [post]
func (controller *MasterController) CreatePort(ctx *gin.Context) {
	var request models.PortRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request pelabuhan tidak valid", err.Error()))
		return
	}

	port, err := controller.masterService.CreatePort(request)
	if err != nil {
		respondMasterError(ctx, "pelabuhan", err)
		return
	}

	ctx.JSON(http.StatusCreated, models.SuccessResponse("pelabuhan berhasil dibuat", port))
}

// GetPort godoc
// @Summary Detail Pelabuhan
// @Description Mengambil satu data master pelabuhan berdasarkan ID (Khusus Officer)
// @Tags Master Pelabuhan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Pelabuhan ID"
// @Success 200 {object} models.APIResponse{data=models.Port} "Pelabuhan Berhasil Diambil"
// @Failure 400 {object} models.APIResponse "ID Tidak Valid"
// @Failure 404 {object} models.APIResponse "Pelabuhan Tidak Ditemukan"
// @Router /master/ports/{id} [get]
func (controller *MasterController) GetPort(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	port, err := controller.masterService.GetPort(id)
	if err != nil {
		respondMasterError(ctx, "pelabuhan", err)
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("pelabuhan berhasil diambil", port))
}

// UpdatePort godoc
// @Summary Perbarui Pelabuhan
// @Description Mengubah data master pelabuhan berdasarkan ID (Khusus Officer)
// @Tags Master Pelabuhan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Pelabuhan ID"
// @Param request body models.PortRequest true "Payload Pelabuhan"
// @Success 200 {object} models.APIResponse{data=models.Port} "Pelabuhan Berhasil Diperbarui"
// @Failure 400 {object} models.APIResponse "Request atau ID Tidak Valid"
// @Failure 404 {object} models.APIResponse "Pelabuhan Tidak Ditemukan"
// @Router /master/ports/{id} [put]
func (controller *MasterController) UpdatePort(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var request models.PortRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request pelabuhan tidak valid", err.Error()))
		return
	}

	port, err := controller.masterService.UpdatePort(id, request)
	if err != nil {
		respondMasterError(ctx, "pelabuhan", err)
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("pelabuhan berhasil diperbarui", port))
}

// DeletePort godoc
// @Summary Hapus Pelabuhan
// @Description Menghapus data master pelabuhan berdasarkan ID (Khusus Officer)
// @Tags Master Pelabuhan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Pelabuhan ID"
// @Success 200 {object} models.APIResponse "Pelabuhan Berhasil Dihapus"
// @Failure 400 {object} models.APIResponse "ID Tidak Valid"
// @Failure 404 {object} models.APIResponse "Pelabuhan Tidak Ditemukan"
// @Router /master/ports/{id} [delete]
func (controller *MasterController) DeletePort(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	if err := controller.masterService.DeletePort(id); err != nil {
		respondMasterError(ctx, "pelabuhan", err)
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("pelabuhan berhasil dihapus", nil))
}

// ListCommodities godoc
// @Summary Daftar Komoditas
// @Description Mengambil semua data master komoditas (Khusus Officer)
// @Tags Master Komoditas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.Commodity} "Komoditas Berhasil Diambil"
// @Failure 500 {object} models.APIResponse "Gagal Mengambil Data"
// @Router /master/commodities [get]
func (controller *MasterController) ListCommodities(ctx *gin.Context) {
	commodities, err := controller.masterService.ListCommodities()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("data komoditas gagal diambil", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("data komoditas berhasil diambil", commodities))
}

// CreateCommodity godoc
// @Summary Tambah Komoditas
// @Description Membuat data master komoditas baru (Khusus Officer)
// @Tags Master Komoditas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CommodityRequest true "Payload Komoditas"
// @Success 201 {object} models.APIResponse{data=models.Commodity} "Komoditas Berhasil Dibuat"
// @Failure 400 {object} models.APIResponse "Request Tidak Valid atau Gagal Diproses"
// @Router /master/commodities [post]
func (controller *MasterController) CreateCommodity(ctx *gin.Context) {
	var request models.CommodityRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request komoditas tidak valid", err.Error()))
		return
	}

	commodity, err := controller.masterService.CreateCommodity(request)
	if err != nil {
		respondMasterError(ctx, "komoditas", err)
		return
	}

	ctx.JSON(http.StatusCreated, models.SuccessResponse("komoditas berhasil dibuat", commodity))
}

// GetCommodity godoc
// @Summary Detail Komoditas
// @Description Mengambil satu data master komoditas berdasarkan ID (Khusus Officer)
// @Tags Master Komoditas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Komoditas ID"
// @Success 200 {object} models.APIResponse{data=models.Commodity} "Komoditas Berhasil Diambil"
// @Failure 400 {object} models.APIResponse "ID Tidak Valid"
// @Failure 404 {object} models.APIResponse "Komoditas Tidak Ditemukan"
// @Router /master/commodities/{id} [get]
func (controller *MasterController) GetCommodity(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	commodity, err := controller.masterService.GetCommodity(id)
	if err != nil {
		respondMasterError(ctx, "komoditas", err)
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("komoditas berhasil diambil", commodity))
}

// UpdateCommodity godoc
// @Summary Perbarui Komoditas
// @Description Mengubah data master komoditas berdasarkan ID (Khusus Officer)
// @Tags Master Komoditas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Komoditas ID"
// @Param request body models.CommodityRequest true "Payload Komoditas"
// @Success 200 {object} models.APIResponse{data=models.Commodity} "Komoditas Berhasil Diperbarui"
// @Failure 400 {object} models.APIResponse "Request atau ID Tidak Valid"
// @Failure 404 {object} models.APIResponse "Komoditas Tidak Ditemukan"
// @Router /master/commodities/{id} [put]
func (controller *MasterController) UpdateCommodity(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var request models.CommodityRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request komoditas tidak valid", err.Error()))
		return
	}

	commodity, err := controller.masterService.UpdateCommodity(id, request)
	if err != nil {
		respondMasterError(ctx, "komoditas", err)
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("komoditas berhasil diperbarui", commodity))
}

// DeleteCommodity godoc
// @Summary Hapus Komoditas
// @Description Menghapus data master komoditas berdasarkan ID (Khusus Officer)
// @Tags Master Komoditas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Komoditas ID"
// @Success 200 {object} models.APIResponse "Komoditas Berhasil Dihapus"
// @Failure 400 {object} models.APIResponse "ID Tidak Valid"
// @Failure 404 {object} models.APIResponse "Komoditas Tidak Ditemukan"
// @Router /master/commodities/{id} [delete]
func (controller *MasterController) DeleteCommodity(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	if err := controller.masterService.DeleteCommodity(id); err != nil {
		respondMasterError(ctx, "komoditas", err)
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("komoditas berhasil dihapus", nil))
}
