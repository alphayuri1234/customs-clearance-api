package controllers

import (
	"net/http"

	"customs-clearance-api/middleware"
	"customs-clearance-api/models"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
)

type WorkflowController struct {
	workflowService *services.WorkflowService
}

func NewWorkflowController(workflowService *services.WorkflowService) *WorkflowController {
	return &WorkflowController{workflowService: workflowService}
}

// InitWorkflow inisialisasi workflow awal (cek risk profile)
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

// ProcessInspection memproses hasil periksa fisik
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

// ProcessApprove menyetujui dokumen clearance
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

// ProcessRelease menerbitkan SPPB
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

// ProcessGateOut memproses pengeluaran barang
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
