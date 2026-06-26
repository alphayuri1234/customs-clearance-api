package routes

import (
	"customs-clearance-api/controllers"
	"customs-clearance-api/middleware"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
)

func RegisterWorkflowRoutes(router *gin.RouterGroup, workflowService *services.WorkflowService) {
	workflowController := controllers.NewWorkflowController(workflowService)

	clearances := router.Group("/clearances")
	clearances.Use(middleware.AuthMiddleware(), middleware.OfficerOnly())
	{
		clearances.GET("", workflowController.ListClearances)
		clearances.GET("/:id", workflowController.GetClearance)
	}

	workflow := router.Group("/workflow")
	// Semua endpoint workflow dilindungi oleh JWT Auth dan hanya untuk Officer
	workflow.Use(middleware.AuthMiddleware(), middleware.OfficerOnly())
	{
		workflow.POST("/init", workflowController.InitWorkflow)
		workflow.POST("/inspection", workflowController.ProcessInspection)
		workflow.POST("/approve", workflowController.ProcessApprove)
		workflow.POST("/release", workflowController.ProcessRelease)
		workflow.POST("/gate-out", workflowController.ProcessGateOut)
	}
}
