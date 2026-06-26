package services

import (
	"errors"
	"fmt"

	"customs-clearance-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkflowService struct {
	db *gorm.DB
}

func NewWorkflowService(db *gorm.DB) *WorkflowService {
	return &WorkflowService{db: db}
}

// EvaluateRiskProfile mengevaluasi tingkat risiko clearance berdasarkan data komoditas dan nilai barang.
func (service *WorkflowService) EvaluateRiskProfile(tx *gorm.DB, clearance *models.Clearance) (*models.RiskProfile, error) {
	// Ambil data detail komoditas dan port (termasuk negara) untuk evaluasi risiko
	var comp models.Commodity
	if err := tx.First(&comp, clearance.CommodityID).Error; err != nil {
		return nil, errors.New("komoditas tidak ditemukan")
	}
	clearance.Commodity = comp

	var port models.Port
	if err := tx.Preload("Country").First(&port, clearance.PortID).Error; err != nil {
		return nil, errors.New("pelabuhan tidak ditemukan")
	}
	clearance.Port = port

	level := models.RiskLevelLow
	score := 0.0
	var reason string

	// Aturan 1: Nilai barang tinggi (> Rp 50.000.000)
	if clearance.Valuation > 50000000 {
		level = models.RiskLevelHigh
		score += 50
		reason += fmt.Sprintf("Nilai barang tinggi (Rp %.2f > Rp 50.000.000); ", clearance.Valuation)
	}

	// Aturan 2: Tarif bea masuk tinggi (> 15%)
	if clearance.Commodity.ImportDutyRate > 15.0 {
		level = models.RiskLevelHigh
		score += 30
		reason += fmt.Sprintf("Tarif bea masuk tinggi (%.2f%% > 15%%); ", clearance.Commodity.ImportDutyRate)
	}

	// Aturan 3: Negara asal tertentu yang butuh pengawasan (contoh negara berkode "HRX")
	if clearance.Port.Country.Code == "HRX" {
		level = models.RiskLevelHigh
		score += 20
		reason += "Negara asal masuk dalam daftar pengawasan khusus (HRX); "
	}

	if reason == "" {
		reason = "Nilai barang dan tarif bea masuk dalam batas aman (Jalur Hijau)."
	}

	riskProfile := &models.RiskProfile{
		ClearanceID: clearance.ID,
		Level:       level,
		Score:       score,
		Reason:      reason,
	}

	return riskProfile, nil
}

// InitWorkflow melakukan pengecekan risk profile saat clearance baru disubmit.
func (service *WorkflowService) InitWorkflow(clearanceID uint) (*models.Clearance, error) {
	var result *models.Clearance

	err := service.db.Transaction(func(tx *gorm.DB) error {
		var clearance models.Clearance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&clearance, clearanceID).Error; err != nil {
			return errors.New("data clearance tidak ditemukan")
		}

		// Pastikan statusnya masih SUBMITTED
		if clearance.Status != models.StatusSubmitted {
			return errors.New("workflow hanya dapat diinisialisasi untuk clearance berstatus SUBMITTED")
		}

		// Periksa apakah risk profile sudah dibuat sebelumnya
		var count int64
		if err := tx.Model(&models.RiskProfile{}).Where("clearance_id = ?", clearanceID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("risk profile untuk clearance ini sudah terdaftar")
		}

		// Evaluasi Risk Profile
		riskProfile, err := service.EvaluateRiskProfile(tx, &clearance)
		if err != nil {
			return err
		}

		// Simpan Risk Profile
		if err := tx.Create(riskProfile).Error; err != nil {
			return err
		}

		// Transisi status berdasarkan tingkat risiko
		if riskProfile.Level == models.RiskLevelHigh {
			clearance.Status = models.StatusInspection
		} else {
			clearance.Status = models.StatusSubmitted
		}

		if err := tx.Save(&clearance).Error; err != nil {
			return err
		}

		result = &clearance
		result.RiskProfile = riskProfile
		return nil
	})

	return result, err
}

// ProcessInspection memproses pemeriksaan fisik. Hanya valid jika status saat ini adalah INSPECTION.
func (service *WorkflowService) ProcessInspection(clearanceID uint, resultInspection string) (*models.Clearance, error) {
	var result *models.Clearance

	err := service.db.Transaction(func(tx *gorm.DB) error {
		var clearance models.Clearance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("RiskProfile").
			First(&clearance, clearanceID).Error; err != nil {
			return errors.New("data clearance tidak ditemukan")
		}

		if clearance.Status != models.StatusInspection {
			return fmt.Errorf("status saat ini '%s'. Transisi ke hasil pemeriksaan hanya valid untuk status 'INSPECTION'", clearance.Status)
		}

		switch resultInspection {
		case "PASS":
			clearance.Status = models.StatusInspectionPassed
		case "FAIL":
			clearance.Status = models.StatusHold
		default:
			return errors.New("hasil pemeriksaan tidak valid, harus 'PASS' atau 'FAIL'")
		}

		if err := tx.Save(&clearance).Error; err != nil {
			return err
		}

		result = &clearance
		return nil
	})

	return result, err
}

// ProcessApprove menyetujui dokumen clearance.
func (service *WorkflowService) ProcessApprove(clearanceID uint) (*models.Clearance, error) {
	var result *models.Clearance

	err := service.db.Transaction(func(tx *gorm.DB) error {
		var clearance models.Clearance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("RiskProfile").
			First(&clearance, clearanceID).Error; err != nil {
			return errors.New("data clearance tidak ditemukan")
		}

		// Pastikan status terkunci HOLD tidak bisa diproses
		if clearance.Status == models.StatusHold {
			return errors.New("clearance dalam status HOLD (terkunci) karena gagal pemeriksaan fisik")
		}

		// Validasi transisi status berdasarkan Risk Profile
		if clearance.RiskProfile == nil {
			return errors.New("risk profile belum dievaluasi. Harap inisialisasi workflow terlebih dahulu")
		}

		if clearance.RiskProfile.Level == models.RiskLevelHigh {
			if clearance.Status != models.StatusInspectionPassed {
				return fmt.Errorf("untuk Risk Profile HIGH, clearance harus berstatus 'INSPECTION_PASSED' sebelum disetujui (status saat ini: %s)", clearance.Status)
			}
		} else { // RiskLevelLow
			if clearance.Status != models.StatusSubmitted {
				return fmt.Errorf("untuk Risk Profile LOW, clearance harus berstatus 'SUBMITTED' sebelum disetujui (status saat ini: %s)", clearance.Status)
			}
		}

		clearance.Status = models.StatusApproved
		if err := tx.Save(&clearance).Error; err != nil {
			return err
		}

		result = &clearance
		return nil
	})

	return result, err
}

// ProcessRelease menerbitkan SPPB (Surat Persetujuan Pengeluaran Barang).
func (service *WorkflowService) ProcessRelease(clearanceID uint) (*models.Clearance, error) {
	var result *models.Clearance

	err := service.db.Transaction(func(tx *gorm.DB) error {
		var clearance models.Clearance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("RiskProfile").
			First(&clearance, clearanceID).Error; err != nil {
			return errors.New("data clearance tidak ditemukan")
		}

		// Untuk LOW risk, bisa melompati APPROVED langsung ke RELEASED dari SUBMITTED
		if clearance.RiskProfile == nil {
			return errors.New("risk profile belum dievaluasi. Harap inisialisasi workflow terlebih dahulu")
		}

		valid := false
		if clearance.Status == models.StatusApproved {
			valid = true
		} else if clearance.RiskProfile.Level == models.RiskLevelLow && clearance.Status == models.StatusSubmitted {
			valid = true
		}

		if !valid {
			return fmt.Errorf("transisi ke RELEASED tidak valid dari status '%s'", clearance.Status)
		}

		clearance.Status = models.StatusReleased
		if err := tx.Save(&clearance).Error; err != nil {
			return err
		}

		result = &clearance
		return nil
	})

	return result, err
}

// ProcessGateOut melakukan proses pengeluaran barang dari kawasan pabean (tahap akhir).
func (service *WorkflowService) ProcessGateOut(clearanceID uint) (*models.Clearance, error) {
	var result *models.Clearance

	err := service.db.Transaction(func(tx *gorm.DB) error {
		var clearance models.Clearance
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("RiskProfile").
			First(&clearance, clearanceID).Error; err != nil {
			return errors.New("data clearance tidak ditemukan")
		}

		if clearance.Status != models.StatusReleased {
			return fmt.Errorf("transisi ke GATE_OUT hanya valid jika status saat ini 'RELEASED' (status saat ini: %s)", clearance.Status)
		}

		clearance.Status = models.StatusGateOut
		if err := tx.Save(&clearance).Error; err != nil {
			return err
		}

		result = &clearance
		return nil
	})

	return result, err
}
