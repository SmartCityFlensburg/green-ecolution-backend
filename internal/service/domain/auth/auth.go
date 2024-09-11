package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type AuthService struct {
	authRepository storage.AuthRepository
	userRepo       storage.UserRepository
	validator      *validator.Validate
	cfg            *config.IdentityAuthConfig
}

func NewAuthService(repo storage.AuthRepository, userRepo storage.UserRepository, cfg *config.IdentityAuthConfig) service.AuthService {
	return &AuthService{
		validator:      validator.New(),
		authRepository: repo,
		userRepo:       userRepo,
		cfg:            cfg,
	}
}

func (s *AuthService) Ready() bool {
	return s.authRepository != nil
}
