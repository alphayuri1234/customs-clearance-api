package routes

import (
	"net/http"

	"customs-clearance-api/models"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	authService := services.NewAuthService(db)
	masterService := services.NewMasterService(db)
	workflowService := services.NewWorkflowService(db)
	seederService := services.NewSeederService(db, workflowService)
	dashboardService := services.NewDashboardService(db)

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, models.SuccessResponse("service aktif", gin.H{
			"service": "customs-clearance-api",
		}))
	})

	router.GET("/swagger-docs/doc.json", func(ctx *gin.Context) {
		ctx.File("./docs/swagger.json")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger-docs/doc.json")))

	api := router.Group("/api/v1")
	RegisterAuthRoutes(api, authService)
	RegisterMasterRoutes(api, masterService)
	RegisterWorkflowRoutes(api, workflowService)
	RegisterDashboardRoutes(api, seederService, dashboardService)

	return router
}
