package services

import (
	"errors"
	"strings"
	"time"

	"customs-clearance-api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		db: db,
	}
}

func (service *AuthService) Register(request models.RegisterRequest) (models.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(request.Email))

	var count int64
	err := service.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return models.AuthResponse{}, err
	}
	if count > 0 {
		return models.AuthResponse{}, errors.New("email sudah terdaftar")
	}

	role := request.Role
	if role == "" {
		role = models.RoleUser
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.AuthResponse{}, err
	}

	now := time.Now()
	user := models.User{
		Name:      strings.TrimSpace(request.Name),
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := service.db.Create(&user).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return models.AuthResponse{}, errors.New("email sudah terdaftar")
		}
		return models.AuthResponse{}, err
	}

	token, err := GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{Token: token, User: user}, nil
}

func (service *AuthService) Login(request models.LoginRequest) (models.AuthResponse, error) {
	var user models.User
	email := strings.ToLower(strings.TrimSpace(request.Email))

	err := service.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.AuthResponse{}, errors.New("email atau password salah")
		}
		return models.AuthResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return models.AuthResponse{}, errors.New("email atau password salah")
	}

	token, err := GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{Token: token, User: user}, nil
}

func (service *AuthService) FindByID(userID uint) (models.User, bool) {
	var user models.User
	err := service.db.First(&user, userID).Error
	if err != nil {
		return models.User{}, false
	}
	return user, true
}
