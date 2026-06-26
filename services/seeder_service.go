package services

import (
	"fmt"
	"math/rand"
	"time"

	"customs-clearance-api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SeederService struct {
	db              *gorm.DB
	workflowService *WorkflowService
}

func NewSeederService(db *gorm.DB, workflowService *WorkflowService) *SeederService {
	return &SeederService{
		db:              db,
		workflowService: workflowService,
	}
}

func (s *SeederService) Seed() error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Bersihkan seluruh data (Truncate dengan CASCADE agar mereset ID serial)
		tables := []string{"risk_profiles", "clearances", "commodities", "ports", "countries", "officers", "users"}
		for _, table := range tables {
			if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
				return fmt.Errorf("gagal mereset tabel %s: %w", table, err)
			}
		}

		// 2. Buat Password Hash
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// 3. Seed Users
		users := []models.User{
			{Name: "PT. Maju Mundur", Email: "importer1@email.com", Password: string(hashedPassword), Role: models.RoleUser},
			{Name: "PT. Elektronik Nusantara", Email: "importer2@email.com", Password: string(hashedPassword), Role: models.RoleUser},
			{Name: "PT. Global Logistik", Email: "importer3@email.com", Password: string(hashedPassword), Role: models.RoleUser},
			{Name: "Supardi (Officer)", Email: "supardi@customs.go.id", Password: string(hashedPassword), Role: models.RoleOfficer},
			{Name: "Hartono (Officer)", Email: "hartono@customs.go.id", Password: string(hashedPassword), Role: models.RoleOfficer},
		}

		for i := range users {
			if err := tx.Create(&users[i]).Error; err != nil {
				return fmt.Errorf("gagal seeding user: %w", err)
			}
		}

		// 4. Seed Officers
		officers := []models.Officer{
			{UserID: users[3].ID, NIP: "199008122015011002", Position: "Senior Inspector"},
			{UserID: users[4].ID, NIP: "198804152010031001", Position: "Document Approver Officer"},
		}
		for i := range officers {
			if err := tx.Create(&officers[i]).Error; err != nil {
				return fmt.Errorf("gagal seeding officer: %w", err)
			}
		}

		// 5. Seed Countries
		countries := []models.Country{
			{Code: "IDN", Name: "Indonesia"},
			{Code: "SGP", Name: "Singapore"},
			{Code: "CHN", Name: "China"},
			{Code: "JPN", Name: "Japan"},
			{Code: "USA", Name: "United States"},
			{Code: "HRX", Name: "High Risk Special Country"}, // Jalur merah pasti terpicu jika dari sini
		}
		for i := range countries {
			if err := tx.Create(&countries[i]).Error; err != nil {
				return fmt.Errorf("gagal seeding country: %w", err)
			}
		}

		// 6. Seed Ports
		ports := []models.Port{
			{Code: "IDJKT", Name: "Tanjung Priok, Jakarta", CountryID: countries[0].ID},
			{Code: "IDSUB", Name: "Tanjung Perak, Surabaya", CountryID: countries[0].ID},
			{Code: "SGPIN", Name: "Jurong Port", CountryID: countries[1].ID},
			{Code: "CNSHA", Name: "Shanghai Port", CountryID: countries[2].ID},
			{Code: "JPTYO", Name: "Tokyo Port", CountryID: countries[3].ID},
			{Code: "USLAX", Name: "Los Angeles Port", CountryID: countries[4].ID},
			{Code: "HRXPT", Name: "Red Port Authority", CountryID: countries[5].ID},
		}
		for i := range ports {
			if err := tx.Create(&ports[i]).Error; err != nil {
				return fmt.Errorf("gagal seeding port: %w", err)
			}
		}

		// 7. Seed Commodities
		commodities := []models.Commodity{
			{HSCode: "85171200", Description: "Telepon Seluler / Smartphones", ImportDutyRate: 10.0, VATRate: 11.0},
			{HSCode: "84713010", Description: "Laptop / Notebook", ImportDutyRate: 0.0, VATRate: 11.0},
			{HSCode: "61091000", Description: "Kaos T-Shirt Katun", ImportDutyRate: 25.0, VATRate: 11.0}, // Tarif Tinggi -> High Risk
			{HSCode: "87032363", Description: "Mobil Sedan Listrik Mewah", ImportDutyRate: 50.0, VATRate: 11.0}, // Tarif Sangat Tinggi -> High Risk
			{HSCode: "90189000", Description: "Alat Pacu Jantung Medis", ImportDutyRate: 2.0, VATRate: 5.0},
			{HSCode: "10063099", Description: "Beras Jasmine Wangi", ImportDutyRate: 5.0, VATRate: 0.0},
		}
		for i := range commodities {
			if err := tx.Create(&commodities[i]).Error; err != nil {
				return fmt.Errorf("gagal seeding commodity: %w", err)
			}
		}

		// 8. Seed Clearances & Trigger Workflow
		// Kita akan membuat 30 data clearance secara acak namun terstruktur agar merata statistiknya.
		descriptions := []string{
			"Importasi Suku Cadang Elektronik", "Pengiriman Pakaian Rajut Musim Dingin",
			"Bahan Baku Medis untuk Rumah Sakit", "Komponen Komputer Industri",
			"Beras untuk Cadangan Pangan Mandiri", "Impor Kendaraan Operasional",
			"Pengiriman Tablet PC untuk Sekolah", "Impor Kain Gulung Katun",
			"Importasi Gadget dan Smartwatch", "Alat Monitoring Bedah Jantung",
		}

		// List importir yang valid (bukan officer)
		importerIDs := []uint{users[0].ID, users[1].ID, users[2].ID}

		// Gunakan seed waktu tetap agar data seeder konsisten setiap dieksekusi
		r := rand.New(rand.NewSource(42))

		// Variasi status untuk didistribusikan
		targetStatuses := []string{
			"SUBMITTED", "INSPECTION", "INSPECTION_PASSED", "APPROVED", "RELEASED", "HOLD", "GATE_OUT",
		}

		for i := 1; i <= 35; i++ {
			importerID := importerIDs[r.Intn(len(importerIDs))]
			commodity := commodities[r.Intn(len(commodities))]
			port := ports[r.Intn(len(ports))]
			desc := fmt.Sprintf("%s Kloter %d", descriptions[r.Intn(len(descriptions))], i)

			// Variasi Valuation: 30% bernilai rendah, 70% bernilai tinggi
			var valuation float64
			if r.Float32() < 0.3 {
				valuation = float64(1000000 + r.Intn(40000000)) // Rp 1 Juta - Rp 41 Juta (Jalur Hijau)
			} else {
				valuation = float64(60000000 + r.Intn(940000000)) // Rp 60 Juta - Rp 1 Miliar (Jalur Merah)
			}

			// Buat clearance dengan status awal SUBMITTED
			createdAt := time.Now().AddDate(0, 0, -r.Intn(10)).Add(time.Duration(-r.Intn(24)) * time.Hour)
			clearance := models.Clearance{
				UserID:      importerID,
				CommodityID: commodity.ID,
				PortID:      port.ID,
				Valuation:   valuation,
				Description: desc,
				Status:      models.StatusSubmitted,
				CreatedAt:   createdAt,
				UpdatedAt:   createdAt,
			}

			// Simpan Clearance mentah
			if err := tx.Create(&clearance).Error; err != nil {
				return fmt.Errorf("gagal membuat mock clearance ke-%d: %w", i, err)
			}

			// Evaluasi Risk Profile & Update Status Awal
			riskProfile, err := s.workflowService.EvaluateRiskProfile(tx, &clearance)
			if err != nil {
				return fmt.Errorf("gagal mengevaluasi risiko clearance ke-%d: %w", i, err)
			}
			if err := tx.Create(riskProfile).Error; err != nil {
				return fmt.Errorf("gagal membuat risk profile clearance ke-%d: %w", i, err)
			}

			if riskProfile.Level == models.RiskLevelHigh {
				clearance.Status = models.StatusInspection
			} else {
				clearance.Status = models.StatusSubmitted
			}
			if err := tx.Save(&clearance).Error; err != nil {
				return err
			}

			// Simulasikan transisi status acak agar sebaran dashboard bervariasi
			targetStatus := targetStatuses[r.Intn(len(targetStatuses))]

			if riskProfile.Level == models.RiskLevelHigh {
				// Skenario HIGH RISK: SUBMITTED -> INSPECTION -> (PASS/FAIL) -> (APPROVED/HOLD) -> RELEASED -> GATE_OUT
				switch targetStatus {
				case "INSPECTION_PASSED":
					clearance.Status = models.StatusInspectionPassed
				case "HOLD":
					clearance.Status = models.StatusHold
				case "APPROVED":
					clearance.Status = models.StatusApproved
				case "RELEASED":
					clearance.Status = models.StatusReleased
				case "GATE_OUT":
					clearance.Status = models.StatusGateOut
				}
			} else {
				// Skenario LOW RISK: SUBMITTED -> APPROVED -> RELEASED -> GATE_OUT
				// Catatan: Low risk melompati INSPECTION & INSPECTION_PASSED & HOLD
				switch targetStatus {
				case "APPROVED":
					clearance.Status = models.StatusApproved
				case "RELEASED", "INSPECTION_PASSED": // mapping dummy
					clearance.Status = models.StatusReleased
				case "GATE_OUT":
					clearance.Status = models.StatusGateOut
				case "HOLD":
					// low risk tidak harus hold, biarkan tetap submitted
					clearance.Status = models.StatusSubmitted
				}
			}

			clearance.UpdatedAt = time.Now()
			if err := tx.Save(&clearance).Error; err != nil {
				return fmt.Errorf("gagal mengupdate status akhir clearance ke-%d: %w", i, err)
			}
		}

		return nil
	})
}
