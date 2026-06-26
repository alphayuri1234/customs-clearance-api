package routes

import (
	"net/http"

	"customs-clearance-api/models"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	authService := services.NewAuthService()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, models.SuccessResponse("service aktif", gin.H{
			"service": "customs-clearance-api",
		}))
	})

	api := router.Group("/api/v1")
	RegisterAuthRoutes(api, authService)

	return router
}
