package controllers

import (
	"net/http"

	"customs-clearance-api/middleware"
	"customs-clearance-api/models"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (controller *AuthController) Register(ctx *gin.Context) {
	var request models.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request register tidak valid", err.Error()))
		return
	}

	response, err := controller.authService.Register(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("registrasi gagal", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, models.SuccessResponse("registrasi berhasil", response))
}

func (controller *AuthController) Login(ctx *gin.Context) {
	var request models.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse("request login tidak valid", err.Error()))
		return
	}

	response, err := controller.authService.Login(request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse("login gagal", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("login berhasil", response))
}

func (controller *AuthController) Me(ctx *gin.Context) {
	claims, ok := middleware.CurrentUserClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse("user belum terautentikasi", nil))
		return
	}

	user, exists := controller.authService.FindByID(claims.UserID)
	if !exists {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse("user tidak ditemukan", nil))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse("profil berhasil diambil", user))
}
