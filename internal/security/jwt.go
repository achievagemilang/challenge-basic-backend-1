//go:generate mockgen -source=jwt.go -destination=../../test/mocks/security_mocks.go -package=mocks
package security

import (
	"fmt"
	"time"

	"challenge-backend-1/internal/entity"

	"github.com/golang-jwt/jwt/v5"
)

type TokenProvider interface {
	GenerateAccessToken(user *entity.User) (string, error)
	GenerateRefreshToken(user *entity.User) (string, error)
	ValidateToken(tokenString string) (*jwt.MapClaims, error)
}

type JwtTokenProvider struct {
	secret string
}

func NewJwtTokenProvider(secret string) *JwtTokenProvider {
	return &JwtTokenProvider{
		secret: secret,
	}
}

func (p *JwtTokenProvider) GenerateAccessToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"type":  "access",
		"exp":   time.Now().Add(time.Second * 20).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(p.secret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (p *JwtTokenProvider) GenerateRefreshToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"type":  "refresh",
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(p.secret))
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (p *JwtTokenProvider) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
