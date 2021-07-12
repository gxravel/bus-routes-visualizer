package jwt

import (
	"context"
	"fmt"

	"github.com/gxravel/bus-routes-visualizer/internal/config"
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	"github.com/gxravel/bus-routes-visualizer/internal/model"
	"github.com/gxravel/bus-routes-visualizer/internal/storage"

	"github.com/dgrijalva/jwt-go"
)

// Manager includes the methods allowed to deal with the token.
type Manager interface {
	parse(tokenString string) (*Claims, error)
	checkIfExist(ctx context.Context, tokenUUID string) error

	Verify(ctx context.Context, tokenString string) (*User, error)
}

// JWT contains the fields which interact with the token.
type JWT struct {
	client *storage.Client
	config config.Config
}

func New(client *storage.Client, config config.Config) *JWT {
	return &JWT{client: client, config: config}
}

// User describes user built into the token
type User struct {
	ID   int64          `json:"id"`
	Type model.UserType `json:"type"`
}

// Claims defines JWT token claims.
type Claims struct {
	User *User `json:"user"`
	jwt.StandardClaims
}

// parse parses a string token with the key and returns claims.
func (m *JWT) parse(tokenString string) (*Claims, error) {
	var key = []byte(m.config.JWT.AccessKey)

	claims := &Claims{}

	jwtToken, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ierr.NewReason(ierr.ErrInvalidJWT).
				WithMessage(fmt.Sprintf("unexpected signing method: %v", t.Header["alg"]))
		}

		return key, nil
	})

	if err != nil || !jwtToken.Valid {
		return nil, ierr.NewReason(ierr.ErrInvalidToken).WithMessage("token validation failed")
	}

	if claims.User == nil {
		return nil, ierr.NewReason(ierr.ErrInvalidToken).WithMessage("failed to get claims user")
	}

	if claims.Id == "" {
		return nil, ierr.NewReason(ierr.ErrInvalidToken).WithMessage("failed to get claims id")
	}

	return claims, nil
}

// checkIfExist checks if token exists in the storage database, and returns an error if not so.
func (m *JWT) checkIfExist(ctx context.Context, tokenUUID string) error {
	return m.client.Get(ctx, tokenUUID).Err()
}

// Verify verifies token, and if it presents in storage returns the user.
func (m *JWT) Verify(ctx context.Context, tokenString string) (*User, error) {
	claims, err := m.parse(tokenString)
	if err != nil {
		return nil, err
	}

	if err := m.checkIfExist(ctx, claims.Id); err != nil {
		return nil, ierr.NewReason(ierr.ErrTokenExpired)
	}

	return claims.User, nil
}
