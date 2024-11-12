package keycloak

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestKeyCloakRepo_RetrospectToken(t *testing.T) {
	t.Run("should return valid token", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		k := NewKeycloakRepository(cfg)

		validToken, err := loginRestAPIClient(ctx, cfg.KeyCloak.BaseURL, cfg.KeyCloak.ClientID, cfg.KeyCloak.ClientSecret, cfg.KeyCloak.Realm)
		assert.NoError(t, err)

		// when
		got, err := k.RetrospectToken(ctx, validToken.AccessToken)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotNil(t, got.Active)
		assert.True(t, *got.Active)
	})

	t.Run("should return not active token when token is invalid", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		k := NewKeycloakRepository(cfg)

		// when
		got, err := k.RetrospectToken(ctx, "invalid-token")

		// then
		assert.NotNil(t, got)
		assert.NotNil(t, got.Active)
		assert.False(t, *got.Active)
		assert.NoError(t, err)
	})

	t.Run("should return not active token when token is empty", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		k := NewKeycloakRepository(cfg)

		// when
		got, err := k.RetrospectToken(ctx, "")

		// then
		assert.NotNil(t, got)
		assert.NotNil(t, got.Active)
		assert.False(t, *got.Active)
		assert.NoError(t, err)
	})
}

func TestKeyCloakRepo_GetAccessTokenFromClientCode(t *testing.T) {
}

func TestKeyCloakRepo_RefreshToken(t *testing.T) {
	t.Run("should return valid token", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		k := NewKeycloakRepository(cfg)
		user := &entities.User{
			Username:    "should-refresh-token",
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       "should-refresh-token@green-ecolution.de",
			EmployeeID:  "123456",
			PhoneNumber: "+49 123456",
		}

		ensureUserExists(t, user)
		validToken := loginUser(t, user)

		// when
		got, err := k.RefreshToken(ctx, validToken.RefreshToken)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotNil(t, got.AccessToken)
		assert.NotNil(t, got.RefreshToken)
		assert.NotEmpty(t, got.AccessToken)
		assert.NotEmpty(t, got.RefreshToken)

		assert.NotEqual(t, validToken.AccessToken, got.AccessToken)
		assert.NotEqual(t, validToken.RefreshToken, got.RefreshToken)

		// validate
		_, err = k.RetrospectToken(ctx, got.AccessToken)
		assert.NoError(t, err)
	})

	t.Run("should return error when refresh token is invalid", func(t *testing.T) {
		// given
		ctx := context.Background()
		cfg := suite.IdentityConfig(t, ctx)
		k := NewKeycloakRepository(cfg)

		// when
		got, err := k.RefreshToken(ctx, "invalid-refresh-token")

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
