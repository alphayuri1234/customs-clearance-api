package services

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dashTestDB *gorm.DB

func TestMainDashboard(t *testing.T) {
	// Diambil dari TestMain yang dipanggil otomatis oleh go test
}

func TestDashboardAndSeeder_Integration(t *testing.T) {
	// Setup database connection (menggunakan testDB dari workflow_service_test.go yang sudah terinisialisasi di TestMain)
	// Namun agar test file ini independen saat dipanggil secara khusus, kita inisialisasi jika nilainya nil.
	db := testDB
	if db == nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("Peringatan: Gagal memuat file .env untuk test dashboard")
		}
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "andiaryatno"
		}
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "customs_clearance"
		}
		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		var dsn string
		if password != "" {
			dsn = "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode + "&TimeZone=Asia/Jakarta"
		} else {
			dsn = "postgres://" + user + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode + "&TimeZone=Asia/Jakarta"
		}

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Gagal terhubung ke database dashboard testing: %v", err)
		}
	}

	workflowService := NewWorkflowService(db)
	seederService := NewSeederService(db, workflowService)
	dashboardService := NewDashboardService(db)

	// 1. Eksekusi Seeder
	err := seederService.Seed()
	if err != nil {
		t.Fatalf("SeederService.Seed() gagal: %v", err)
	}

	// 2. Ambil summary dashboard
	summary, err := dashboardService.GetSummary()
	if err != nil {
		t.Fatalf("DashboardService.GetSummary() gagal: %v", err)
	}

	// 3. Verifikasi jumlah data
	if summary.TotalClearances != 35 {
		t.Errorf("Total clearance harus 35, got: %d", summary.TotalClearances)
	}

	if summary.TotalValuation <= 0 {
		t.Errorf("Total valuation harus lebih besar dari 0, got: %f", summary.TotalValuation)
	}

	// 4. Verifikasi distribusi status
	var totalStatusSum int64
	for status, count := range summary.StatusCounts {
		totalStatusSum += count
		log.Printf("Status %s: %d", status, count)
	}
	if totalStatusSum != 35 {
		t.Errorf("Jumlah data di status counts harus 35, got: %d", totalStatusSum)
	}

	// 5. Verifikasi kategori risiko
	var totalRiskSum int64
	for risk, count := range summary.RiskCounts {
		totalRiskSum += count
		log.Printf("Risk %s: %d", risk, count)
	}
	if totalRiskSum != 35 {
		t.Errorf("Jumlah data di risk counts harus 35, got: %d", totalRiskSum)
	}

	// 6. Verifikasi Top Commodities & Ports
	if len(summary.TopCommodities) == 0 {
		t.Error("Top Commodities tidak boleh kosong")
	}
	if len(summary.TopPorts) == 0 {
		t.Error("Top Ports tidak boleh kosong")
	}

	// 7. Verifikasi Recent Clearances
	if len(summary.RecentClearances) != 5 {
		t.Errorf("Recent clearances harus berisi tepat 5 data, got: %d", len(summary.RecentClearances))
	}

	// Cek apakah data recent clearance ter-preload dengan benar
	firstClearance := summary.RecentClearances[0]
	if firstClearance.User.ID == 0 || firstClearance.Commodity.ID == 0 || firstClearance.Port.ID == 0 || firstClearance.RiskProfile == nil {
		t.Error("Relasi (User, Commodity, Port, RiskProfile) pada recent clearance harus di-preload")
	}
}
