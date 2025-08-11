package jwt

import (
	"time"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenService struct {
	secretKey []byte
}

func NewJWTTokenService(secretKey string) domain.TokenService {
	return &JWTTokenService{
		secretKey: []byte(secretKey),
	}
}

func (t *JWTTokenService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 horas
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.secretKey)
}

func (t *JWTTokenService) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return t.secretKey, nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}