package routes

import (
	"customs-clearance-api/controllers"
	"customs-clearance-api/middleware"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
)

func RegisterDashboardRoutes(router *gin.RouterGroup, seederService *services.SeederService, dashboardService *services.DashboardService) {
	dashboardController := controllers.NewDashboardController(seederService, dashboardService)

	// Route seed dibiarkan public demi kemudahan testing/setup awal di local
	router.POST("/seed", dashboardController.SeedData)

	// Route dashboard diproteksi khusus Officer
	router.GET("/dashboard", middleware.AuthMiddleware(), middleware.OfficerOnly(), dashboardController.GetSummary)
}
