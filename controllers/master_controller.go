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

func (controller *MasterController) ListCountries(ctx *gin.Context) {
	countries, err := controller.masterService.ListCountries()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("data negara gagal diambil", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("data negara berhasil diambil", countries))
}

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

func (controller *MasterController) ListPorts(ctx *gin.Context) {
	ports, err := controller.masterService.ListPorts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("data pelabuhan gagal diambil", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("data pelabuhan berhasil diambil", ports))
}

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

func (controller *MasterController) ListCommodities(ctx *gin.Context) {
	commodities, err := controller.masterService.ListCommodities()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("data komoditas gagal diambil", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("data komoditas berhasil diambil", commodities))
}

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
