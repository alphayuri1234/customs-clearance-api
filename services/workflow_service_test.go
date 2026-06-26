package services

import (
	"fmt"
	"log"
	"os"
	"testing"

	"customs-clearance-api/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	// Muat file .env dari parent directory
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Peringatan: Gagal memuat file .env untuk test, menggunakan env default")
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
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=Asia/Jakarta",
			user, password, host, port, dbname, sslmode,
		)
	} else {
		dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s&TimeZone=Asia/Jakarta",
			user, host, port, dbname, sslmode,
		)
	}

	var err error
	log.Printf("DEBUG DSN: %s", dsn)
	testDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database testing: %v", err)
	}

	var dbName string
	testDB.Raw("SELECT current_database()").Scan(&dbName)
	log.Printf("DEBUG CONNECTED DB: %s", dbName)

	// Migrasikan database
	err = testDB.AutoMigrate(
		&models.User{},
		&models.Officer{},
		&models.Country{},
		&models.Port{},
		&models.Commodity{},
		&models.Clearance{},
		&models.RiskProfile{},
	)
	if err != nil {
		log.Fatalf("Gagal AutoMigrate di testing: %v", err)
	}

	os.Exit(m.Run())
}

func cleanupDatabase() {
	testDB.Exec("DELETE FROM risk_profiles")
	testDB.Exec("DELETE FROM clearances")
	testDB.Exec("DELETE FROM commodities")
	testDB.Exec("DELETE FROM ports")
	testDB.Exec("DELETE FROM countries")
	testDB.Exec("DELETE FROM officers")
	testDB.Exec("DELETE FROM users")
}

func setupTestData(t *testing.T) (models.User, models.Country, models.Port, models.Commodity) {
	cleanupDatabase()

	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     models.RoleUser,
	}
	if err := testDB.Create(&user).Error; err != nil {
		t.Fatalf("Gagal membuat test user: %v", err)
	}

	country := models.Country{
		Code: "IDN",
		Name: "Indonesia",
	}
	if err := testDB.Create(&country).Error; err != nil {
		t.Fatalf("Gagal membuat test country: %v", err)
	}

	port := models.Port{
		Code:      "IDJKT",
		Name:      "Tanjung Priok",
		CountryID: country.ID,
	}
	if err := testDB.Create(&port).Error; err != nil {
		t.Fatalf("Gagal membuat test port: %v", err)
	}

	commodity := models.Commodity{
		HSCode:         "85171200",
		Description:    "Telepon Seluler",
		ImportDutyRate: 10.0,
		VATRate:        11.0,
	}
	if err := testDB.Create(&commodity).Error; err != nil {
		t.Fatalf("Gagal membuat test commodity: %v", err)
	}

	return user, country, port, commodity
}

func TestWorkflowHighRisk_Success(t *testing.T) {
	user, _, port, commodity := setupTestData(t)
	workflowService := NewWorkflowService(testDB)

	// Buat Clearance bernilai tinggi (HIGH Risk)
	clearance := models.Clearance{
		UserID:      user.ID,
		CommodityID: commodity.ID,
		PortID:      port.ID,
		Valuation:   60000000.0, // Rp 60 juta (> 50 juta)
		Description: "Barang Impor Elektronik",
		Status:      models.StatusSubmitted,
	}
	if err := testDB.Create(&clearance).Error; err != nil {
		t.Fatalf("Gagal membuat clearance: %v", err)
	}

	// 1. Inisialisasi Workflow
	initCl, err := workflowService.InitWorkflow(clearance.ID)
	if err != nil {
		t.Fatalf("InitWorkflow gagal: %v", err)
	}

	if initCl.Status != models.StatusInspection {
		t.Errorf("Harus berstatus INSPECTION untuk HIGH risk, got: %s", initCl.Status)
	}

	if initCl.RiskProfile == nil || initCl.RiskProfile.Level != models.RiskLevelHigh {
		t.Errorf("RiskProfile level harus HIGH, got: %v", initCl.RiskProfile)
	}

	// 2. Coba approve langsung -> Harus error karena status masih INSPECTION
	_, err = workflowService.ProcessApprove(clearance.ID)
	if err == nil {
		t.Error("Harus error jika approve dipanggil langsung untuk HIGH risk tanpa inspection")
	}

	// 3. Process Inspection PASS -> Status INSPECTION_PASSED
	inspCl, err := workflowService.ProcessInspection(clearance.ID, "PASS")
	if err != nil {
		t.Fatalf("ProcessInspection gagal: %v", err)
	}
	if inspCl.Status != models.StatusInspectionPassed {
		t.Errorf("Status harus INSPECTION_PASSED, got: %s", inspCl.Status)
	}

	// 4. Process Approve -> Status APPROVED
	appCl, err := workflowService.ProcessApprove(clearance.ID)
	if err != nil {
		t.Fatalf("ProcessApprove gagal: %v", err)
	}
	if appCl.Status != models.StatusApproved {
		t.Errorf("Status harus APPROVED, got: %s", appCl.Status)
	}

	// 5. Process Release -> Status RELEASED
	relCl, err := workflowService.ProcessRelease(clearance.ID)
	if err != nil {
		t.Fatalf("ProcessRelease gagal: %v", err)
	}
	if relCl.Status != models.StatusReleased {
		t.Errorf("Status harus RELEASED, got: %s", relCl.Status)
	}

	// 6. Process Gate Out -> Status GATE_OUT
	goCl, err := workflowService.ProcessGateOut(clearance.ID)
	if err != nil {
		t.Fatalf("ProcessGateOut gagal: %v", err)
	}
	if goCl.Status != models.StatusGateOut {
		t.Errorf("Status harus GATE_OUT, got: %s", goCl.Status)
	}
}

func TestWorkflowLowRisk_Success(t *testing.T) {
	user, _, port, commodity := setupTestData(t)
	workflowService := NewWorkflowService(testDB)

	// Buat Clearance bernilai rendah (LOW Risk)
	clearance := models.Clearance{
		UserID:      user.ID,
		CommodityID: commodity.ID,
		PortID:      port.ID,
		Valuation:   10000000.0, // Rp 10 juta (<= 50 juta)
		Description: "Barang Impor Buku",
		Status:      models.StatusSubmitted,
	}
	if err := testDB.Create(&clearance).Error; err != nil {
		t.Fatalf("Gagal membuat clearance: %v", err)
	}

	// 1. Inisialisasi Workflow
	initCl, err := workflowService.InitWorkflow(clearance.ID)
	if err != nil {
		t.Fatalf("InitWorkflow gagal: %v", err)
	}

	if initCl.Status != models.StatusSubmitted {
		t.Errorf("Harus tetap berstatus SUBMITTED untuk LOW risk, got: %s", initCl.Status)
	}

	if initCl.RiskProfile == nil || initCl.RiskProfile.Level != models.RiskLevelLow {
		t.Errorf("RiskProfile level harus LOW, got: %v", initCl.RiskProfile)
	}

	// 2. Coba periksa fisik -> Harus error karena LOW risk tidak masuk antrean INSPECTION
	_, err = workflowService.ProcessInspection(clearance.ID, "PASS")
	if err == nil {
		t.Error("Harus error jika periksa fisik dipanggil untuk status SUBMITTED")
	}

	// 3. Process Approve langsung -> Status APPROVED (Melewati pemeriksaan fisik)
	appCl, err := workflowService.ProcessApprove(clearance.ID)
	if err != nil {
		t.Fatalf("ProcessApprove langsung gagal: %v", err)
	}
	if appCl.Status != models.StatusApproved {
		t.Errorf("Status harus APPROVED, got: %s", appCl.Status)
	}

	// 4. Process Release -> Status RELEASED
	relCl, err := workflowService.ProcessRelease(clearance.ID)
	if err != nil {
		t.Fatalf("ProcessRelease gagal: %v", err)
	}
	if relCl.Status != models.StatusReleased {
		t.Errorf("Status harus RELEASED, got: %s", relCl.Status)
	}
}

func TestWorkflowHighRisk_FailInspection(t *testing.T) {
	user, _, port, commodity := setupTestData(t)
	workflowService := NewWorkflowService(testDB)

	clearance := models.Clearance{
		UserID:      user.ID,
		CommodityID: commodity.ID,
		PortID:      port.ID,
		Valuation:   75000000.0,
		Description: "Barang Mewah Elektronik",
		Status:      models.StatusSubmitted,
	}
	if err := testDB.Create(&clearance).Error; err != nil {
		t.Fatalf("Gagal membuat clearance: %v", err)
	}

	// 1. Inisialisasi
	_, err := workflowService.InitWorkflow(clearance.ID)
	if err != nil {
		t.Fatalf("InitWorkflow gagal: %v", err)
	}

	// 2. Process Inspection FAIL -> Status HOLD
	inspCl, err := workflowService.ProcessInspection(clearance.ID, "FAIL")
	if err != nil {
		t.Fatalf("ProcessInspection gagal: %v", err)
	}
	if inspCl.Status != models.StatusHold {
		t.Errorf("Status harus HOLD, got: %s", inspCl.Status)
	}

	// 3. Coba Approve -> Harus gagal karena berstatus HOLD
	_, err = workflowService.ProcessApprove(clearance.ID)
	if err == nil {
		t.Error("Harus error jika disetujui dalam kondisi HOLD")
	}
}
