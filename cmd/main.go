package main

import (
	"log"
	"os"

	"customs-clearance-api/database"
	"customs-clearance-api/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Gagal memuat file .env, menggunakan env default sistem")
	}

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
