package routes

import (
	"customs-clearance-api/controllers"
	"customs-clearance-api/middleware"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, authService *services.AuthService) {
	authController := controllers.NewAuthController(authService)

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.GET("/me", middleware.AuthMiddleware(), authController.Me)
}
