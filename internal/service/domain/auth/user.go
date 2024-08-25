package auth

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/pkg/errors"
)

func (s *AuthService) Register(ctx context.Context, user *auth.RegisterUser) (*auth.User, error) {
	if err := s.validator.Struct(user); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	createdUser, err := s.authRepository.CreateUser(ctx, &user.User, user.Password, user.Roles)
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to create user").Error())
	}

	return createdUser, nil
}
