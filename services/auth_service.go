package services

import (
	"errors"
	"strings"
	"sync"
	"time"

	"customs-clearance-api/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	mu     sync.RWMutex
	users  map[string]models.User
	nextID uint
}

func NewAuthService() *AuthService {
	return &AuthService{
		users:  make(map[string]models.User),
		nextID: 1,
	}
}

func (service *AuthService) Register(request models.RegisterRequest) (models.AuthResponse, error) {
	service.mu.Lock()
	defer service.mu.Unlock()

	email := strings.ToLower(strings.TrimSpace(request.Email))
	if _, exists := service.users[email]; exists {
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
		ID:        service.nextID,
		Name:      strings.TrimSpace(request.Name),
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}

	service.users[email] = user
	service.nextID++

	token, err := GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{Token: token, User: user}, nil
}

func (service *AuthService) Login(request models.LoginRequest) (models.AuthResponse, error) {
	service.mu.RLock()
	user, exists := service.users[strings.ToLower(strings.TrimSpace(request.Email))]
	service.mu.RUnlock()

	if !exists {
		return models.AuthResponse{}, errors.New("email atau password salah")
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
	service.mu.RLock()
	defer service.mu.RUnlock()

	for _, user := range service.users {
		if user.ID == userID {
			return user, true
		}
	}

	return models.User{}, false
}
