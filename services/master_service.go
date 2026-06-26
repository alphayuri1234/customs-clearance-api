package services

import (
	"errors"
	"strings"

	"customs-clearance-api/models"
	"gorm.io/gorm"
)

type MasterService struct {
	db *gorm.DB
}

func NewMasterService(db *gorm.DB) *MasterService {
	return &MasterService{db: db}
}

func normalizeCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func normalizeText(value string) string {
	return strings.TrimSpace(value)
}

func duplicateError(err error, message string) error {
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		return errors.New(message)
	}
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique") {
		return errors.New(message)
	}
	return err
}

func (service *MasterService) ListCountries() ([]models.Country, error) {
	var countries []models.Country
	err := service.db.Order("name ASC").Find(&countries).Error
	return countries, err
}

func (service *MasterService) CreateCountry(request models.CountryRequest) (models.Country, error) {
	country := models.Country{
		Code: normalizeCode(request.Code),
		Name: normalizeText(request.Name),
	}

	err := service.db.Create(&country).Error
	return country, duplicateError(err, "kode negara sudah terdaftar")
}

func (service *MasterService) GetCountry(id uint) (models.Country, error) {
	var country models.Country
	err := service.db.First(&country, id).Error
	return country, err
}

func (service *MasterService) UpdateCountry(id uint, request models.CountryRequest) (models.Country, error) {
	country, err := service.GetCountry(id)
	if err != nil {
		return models.Country{}, err
	}

	country.Code = normalizeCode(request.Code)
	country.Name = normalizeText(request.Name)

	err = service.db.Save(&country).Error
	return country, duplicateError(err, "kode negara sudah terdaftar")
}

func (service *MasterService) DeleteCountry(id uint) error {
	err := service.db.Delete(&models.Country{}, id).Error
	return err
}

func (service *MasterService) ListPorts() ([]models.Port, error) {
	var ports []models.Port
	err := service.db.Preload("Country").Order("name ASC").Find(&ports).Error
	return ports, err
}

func (service *MasterService) CreatePort(request models.PortRequest) (models.Port, error) {
	if _, err := service.GetCountry(request.CountryID); err != nil {
		return models.Port{}, errors.New("negara tidak ditemukan")
	}

	port := models.Port{
		Code:      normalizeCode(request.Code),
		Name:      normalizeText(request.Name),
		CountryID: request.CountryID,
	}

	err := service.db.Create(&port).Error
	if err != nil {
		return models.Port{}, duplicateError(err, "kode pelabuhan sudah terdaftar")
	}

	err = service.db.Preload("Country").First(&port, port.ID).Error
	return port, err
}

func (service *MasterService) GetPort(id uint) (models.Port, error) {
	var port models.Port
	err := service.db.Preload("Country").First(&port, id).Error
	return port, err
}

func (service *MasterService) UpdatePort(id uint, request models.PortRequest) (models.Port, error) {
	if _, err := service.GetCountry(request.CountryID); err != nil {
		return models.Port{}, errors.New("negara tidak ditemukan")
	}

	port, err := service.GetPort(id)
	if err != nil {
		return models.Port{}, err
	}

	port.Code = normalizeCode(request.Code)
	port.Name = normalizeText(request.Name)
	port.CountryID = request.CountryID

	err = service.db.Save(&port).Error
	if err != nil {
		return models.Port{}, duplicateError(err, "kode pelabuhan sudah terdaftar")
	}

	err = service.db.Preload("Country").First(&port, port.ID).Error
	return port, err
}

func (service *MasterService) DeletePort(id uint) error {
	err := service.db.Delete(&models.Port{}, id).Error
	return err
}

func (service *MasterService) ListCommodities() ([]models.Commodity, error) {
	var commodities []models.Commodity
	err := service.db.Order("hs_code ASC").Find(&commodities).Error
	return commodities, err
}

func (service *MasterService) CreateCommodity(request models.CommodityRequest) (models.Commodity, error) {
	commodity := models.Commodity{
		HSCode:         normalizeCode(request.HSCode),
		Description:    normalizeText(request.Description),
		ImportDutyRate: request.ImportDutyRate,
		VATRate:        request.VATRate,
	}

	err := service.db.Create(&commodity).Error
	return commodity, duplicateError(err, "kode HS sudah terdaftar")
}

func (service *MasterService) GetCommodity(id uint) (models.Commodity, error) {
	var commodity models.Commodity
	err := service.db.First(&commodity, id).Error
	return commodity, err
}

func (service *MasterService) UpdateCommodity(id uint, request models.CommodityRequest) (models.Commodity, error) {
	commodity, err := service.GetCommodity(id)
	if err != nil {
		return models.Commodity{}, err
	}

	commodity.HSCode = normalizeCode(request.HSCode)
	commodity.Description = normalizeText(request.Description)
	commodity.ImportDutyRate = request.ImportDutyRate
	commodity.VATRate = request.VATRate

	err = service.db.Save(&commodity).Error
	return commodity, duplicateError(err, "kode HS sudah terdaftar")
}

func (service *MasterService) DeleteCommodity(id uint) error {
	err := service.db.Delete(&models.Commodity{}, id).Error
	return err
}
