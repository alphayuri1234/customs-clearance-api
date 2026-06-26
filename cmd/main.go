package main

import (
	"log"
	"os"

	"customs-clearance-api/database"
	"customs-clearance-api/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	db := database.InitDB()

	router := routes.SetupRouter(db)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}
