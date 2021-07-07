package jwt

import (
	"context"
	"fmt"

	"github.com/gxravel/bus-routes-visualizer/internal/config"
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	"github.com/gxravel/bus-routes-visualizer/internal/storage"

	"github.com/dgrijalva/jwt-go"
)

// Manager includes the methods allowed to deal with the token.
type Manager interface {
	parse(tokenString string) (string, error)
	checkIfExist(ctx context.Context, tokenUUID string) error

	Verify(ctx context.Context, tokenString string) (int64, error)
}

// JWT contains the fields which interact with the token.
type JWT struct {
	client *storage.Client
	config config.Config
}

func New(client *storage.Client, config config.Config) *JWT {
	return &JWT{client: client, config: config}
}

// parse parses a string token with the key.
func (m *JWT) parse(tokenString string) (string, error) {
	var key = []byte(m.config.JWT.AccessKey)

	claims := jwt.MapClaims{}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ierr.NewReason(ierr.ErrInvalidJWT).
				WithMessage(fmt.Sprintf("unexpected signing method: %v", t.Header["alg"]))
		}

		return key, nil
	})

	if err != nil || !jwtToken.Valid {
		return "", ierr.NewReason(ierr.ErrInvalidToken).WithMessage("token validation failed")
	}

	tokenUUID, ok := claims["jti"].(string)
	if !ok {
		return "", ierr.NewReason(ierr.ErrInvalidToken).WithMessage("failed to get claims id")
	}

	return tokenUUID, nil
}

// checkIfExist checks if token exists in the storage database and re.
func (m *JWT) checkIfExist(ctx context.Context, tokenUUID string) error {
	return m.client.Get(ctx, tokenUUID).Err()
}

// Verify verifies token, and if it presents in storage returns the user id.
func (m *JWT) Verify(ctx context.Context, tokenString string) (int64, error) {
	tokenUUID, err := m.parse(tokenString)
	if err != nil {
		return 0, err
	}

	userID, err := m.client.Get(ctx, tokenUUID).Int64()
	if err != nil {
		return 0, ierr.NewReason(ierr.ErrTokenExpired)
	}

	return userID, nil
}
