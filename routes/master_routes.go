package routes

import (
	"customs-clearance-api/controllers"
	"customs-clearance-api/middleware"
	"customs-clearance-api/services"
	"github.com/gin-gonic/gin"
)

func RegisterMasterRoutes(router *gin.RouterGroup, masterService *services.MasterService) {
	masterController := controllers.NewMasterController(masterService)

	master := router.Group("/master")
	master.Use(middleware.AuthMiddleware(), middleware.OfficerOnly())
	{
		countries := master.Group("/countries")
		{
			countries.GET("", masterController.ListCountries)
			countries.POST("", masterController.CreateCountry)
			countries.GET("/:id", masterController.GetCountry)
			countries.PUT("/:id", masterController.UpdateCountry)
			countries.DELETE("/:id", masterController.DeleteCountry)
		}

		ports := master.Group("/ports")
		{
			ports.GET("", masterController.ListPorts)
			ports.POST("", masterController.CreatePort)
			ports.GET("/:id", masterController.GetPort)
			ports.PUT("/:id", masterController.UpdatePort)
			ports.DELETE("/:id", masterController.DeletePort)
		}

		commodities := master.Group("/commodities")
		{
			commodities.GET("", masterController.ListCommodities)
			commodities.POST("", masterController.CreateCommodity)
			commodities.GET("/:id", masterController.GetCommodity)
			commodities.PUT("/:id", masterController.UpdateCommodity)
			commodities.DELETE("/:id", masterController.DeleteCommodity)
		}
	}
}
