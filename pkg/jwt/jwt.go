package jwt

import (
	"errors"
	"time"

	"github.com/YelzhanWeb/uno-spicchio/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int         `json:"user_id"`
	Username string      `json:"username"`
	Role     domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret     []byte
	expiration time.Duration
}

func NewTokenManager(secret string, expiration time.Duration) *TokenManager {
	return &TokenManager{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

func (tm *TokenManager) Generate(userID int, username string, role domain.Role) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.secret)
}

func (tm *TokenManager) Verify(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return tm.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
