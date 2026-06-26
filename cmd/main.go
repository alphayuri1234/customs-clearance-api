package main

import (
	"log"
	"os"

	"customs-clearance-api/database"
	"customs-clearance-api/routes"
	"github.com/joho/godotenv"
)

// @title Customs Clearance API
// @version 1.0
// @description API untuk sistem Customs Clearance Bea Cukai (Simulasi Jalur Merah & Jalur Hijau).
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Tipe token: Bearer <JWT_TOKEN>
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
