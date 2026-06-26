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

// Register godoc
// @Summary Registrasi Akun Baru
// @Description Membuat akun importir umum (User) atau petugas bea cukai (Officer) baru
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Payload Registrasi"
// @Success 201 {object} models.APIResponse{data=models.AuthResponse} "Registrasi Berhasil"
// @Failure 400 {object} models.APIResponse "Request Tidak Valid atau Registrasi Gagal"
// @Router /register [post]
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

// Login godoc
// @Summary Login Akun
// @Description Melakukan login untuk mendapatkan token akses JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Payload Login"
// @Success 200 {object} models.APIResponse{data=models.AuthResponse} "Login Berhasil"
// @Failure 400 {object} models.APIResponse "Request Tidak Valid"
// @Failure 401 {object} models.APIResponse "Kredensial Tidak Valid"
// @Router /login [post]
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

// Me godoc
// @Summary Profil Saya (Me)
// @Description Mengambil data profil user/officer yang sedang login berdasarkan token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.User} "Profil Berhasil Diambil"
// @Failure 401 {object} models.APIResponse "User Belum Terautentikasi"
// @Failure 404 {object} models.APIResponse "User Tidak Ditemukan"
// @Router /me [get]
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
