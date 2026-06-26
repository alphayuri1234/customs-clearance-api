package controllers

import (
	"errors"
	"net/http"

	"customs-clearance-api/middleware"
	"customs-clearance-api/models"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowController struct {
	workflowService *services.WorkflowService
}

func NewWorkflowController(workflowService *services.WorkflowService) *WorkflowController {
	return &WorkflowController{workflowService: workflowService}
}

// ListClearances godoc
// @Summary Daftar Customs Clearance
// @Description Menampilkan daftar clearance dengan pagination dan filter (Khusus Officer)
// @Tags Workflow Customs Clearance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Nomor halaman"
// @Param limit query int false "Jumlah data per halaman, maksimal 100"
// @Param status query string false "Filter status clearance"
// @Param user_id query int false "Filter user importir"
// @Param commodity_id query int false "Filter komoditas"
// @Param port_id query int false "Filter pelabuhan"
// @Param risk_level query string false "Filter risk level LOW/HIGH"
// @Param search query string false "Cari deskripsi, HS code, komoditas, atau pelabuhan"
// @Success 200 {object} models.APIResponse{data=models.ClearanceListResponse} "Data Clearance Berhasil Diambil"
// @Failure 400 {object} models.APIResponse "Query Tidak Valid"
// @Failure 500 {object} models.APIResponse "Gagal Mengambil Data"
// @Router /clearances [get]
func (controller *WorkflowController) ListClearances(ctx *gin.Context) {
	var query models.ClearanceListQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("query clearance tidak valid", err.Error()))
		return
	}

	response, err := controller.workflowService.ListClearances(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("data clearance gagal diambil", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("data clearance berhasil diambil", response))
}

// GetClearance godoc
// @Summary Detail Customs Clearance
// @Description Menampilkan detail clearance berdasarkan ID (Khusus Officer)
// @Tags Workflow Customs Clearance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Clearance"
// @Success 200 {object} models.APIResponse{data=models.Clearance} "Detail Clearance Berhasil Diambil"
// @Failure 400 {object} models.APIResponse "ID Tidak Valid"
// @Failure 404 {object} models.APIResponse "Clearance Tidak Ditemukan"
// @Router /clearances/{id} [get]
func (controller *WorkflowController) GetClearance(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	clearance, err := controller.workflowService.GetClearance(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse("clearance tidak ditemukan", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse("detail clearance gagal diambil", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("detail clearance berhasil diambil", clearance))
}

// InitWorkflow godoc
// @Summary Inisialisasi Workflow
// @Description Memulai proses evaluasi Risk Engine pada clearance yang baru masuk (SUBMITTED) (Khusus Officer)
// @Tags Workflow Customs Clearance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object{clearance_id=int} true "Payload Inisialisasi Workflow"
// @Success 200 {object} models.APIResponse "Workflow Berhasil Diinisialisasi"
// @Failure 400 {object} models.APIResponse "Request Tidak Valid atau Inisialisasi Gagal"
// @Router /workflow/init [post]
func (controller *WorkflowController) InitWorkflow(ctx *gin.Context) {
	var request struct {
		ClearanceID uint `json:"clearance_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request tidak valid", err.Error()))
		return
	}

	clearance, err := controller.workflowService.InitWorkflow(request.ClearanceID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("inisialisasi workflow gagal", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("workflow berhasil diinisialisasi", gin.H{
		"id":           clearance.ID,
		"status":       clearance.Status,
		"risk_profile": clearance.RiskProfile,
	}))
}

// ProcessInspection godoc
// @Summary Proses Hasil Pemeriksaan Fisik
// @Description Menginputkan hasil periksa fisik (PASS/FAIL) untuk clearance berisiko tinggi (INSPECTION) (Khusus Officer)
// @Tags Workflow Customs Clearance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.InspectionRequest true "Payload Hasil Pemeriksaan"
// @Success 200 {object} models.APIResponse "Hasil Pemeriksaan Berhasil Diproses"
// @Failure 400 {object} models.APIResponse "Request Tidak Valid atau Gagal Diproses"
// @Router /workflow/inspection [post]
func (controller *WorkflowController) ProcessInspection(ctx *gin.Context) {
	var request models.InspectionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request pemeriksaan tidak valid", err.Error()))
		return
	}

	claims, ok := middleware.CurrentUserClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse("user belum terautentikasi", nil))
		return
	}

	clearance, err := controller.workflowService.ProcessInspection(request.ClearanceID, request.Result, claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("pemeriksaan fisik gagal diproses", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("hasil periksa fisik berhasil diproses", gin.H{
		"id":     clearance.ID,
		"status": clearance.Status,
	}))
}

// ProcessApprove godoc
// @Summary Persetujuan Customs Clearance
// @Description Memberikan persetujuan Bea Cukai pada dokumen clearance (Khusus Officer)
// @Tags Workflow Customs Clearance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ApproveRequest true "Payload Persetujuan"
// @Success 200 {object} models.APIResponse "Clearance Berhasil Disetujui"
// @Failure 400 {object} models.APIResponse "Request Tidak Valid atau Gagal Diproses"
// @Router /workflow/approve [post]
func (controller *WorkflowController) ProcessApprove(ctx *gin.Context) {
	var request models.ApproveRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request approval tidak valid", err.Error()))
		return
	}

	clearance, err := controller.workflowService.ProcessApprove(request.ClearanceID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("gagal menyetujui clearance", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("clearance berhasil disetujui", gin.H{
		"id":     clearance.ID,
		"status": clearance.Status,
	}))
}

// ProcessRelease godoc
// @Summary Menerbitkan SPPB (Surat Persetujuan Pengeluaran Barang)
// @Description Mengeluarkan SPPB untuk kontainer agar bisa keluar dari pelabuhan (Khusus Officer)
// @Tags Workflow Customs Clearance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ReleaseRequest true "Payload Rilis SPPB"
// @Success 200 {object} models.APIResponse "SPPB Berhasil Diterbitkan"
// @Failure 400 {object} models.APIResponse "Request Tidak Valid atau Gagal Diproses"
// @Router /workflow/release [post]
func (controller *WorkflowController) ProcessRelease(ctx *gin.Context) {
	var request models.ReleaseRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request release tidak valid", err.Error()))
		return
	}

	claims, ok := middleware.CurrentUserClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse("user belum terautentikasi", nil))
		return
	}

	clearance, err := controller.workflowService.ProcessRelease(request.ClearanceID, claims.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("gagal menerbitkan SPPB", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("SPPB berhasil diterbitkan", gin.H{
		"id":     clearance.ID,
		"status": clearance.Status,
	}))
}

// ProcessGateOut godoc
// @Summary Keluar Gerbang (Gate Out)
// @Description Memproses pengeluaran fisik barang/kontainer dari kawasan pabean pelabuhan (Khusus Officer)
// @Tags Workflow Customs Clearance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.GateOutRequest true "Payload Gate Out"
// @Success 200 {object} models.APIResponse "Barang Berhasil Keluar"
// @Failure 400 {object} models.APIResponse "Request Tidak Valid atau Gagal Diproses"
// @Router /workflow/gate-out [post]
func (controller *WorkflowController) ProcessGateOut(ctx *gin.Context) {
	var request models.GateOutRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request gate out tidak valid", err.Error()))
		return
	}

	clearance, err := controller.workflowService.ProcessGateOut(request.ClearanceID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("gagal memproses gate out", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("barang berhasil keluar dari kawasan pabean", gin.H{
		"id":     clearance.ID,
		"status": clearance.Status,
	}))
}
