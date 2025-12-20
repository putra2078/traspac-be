package utils

import (
	"errors"
	"time"

	"hrm-app/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateTokens(cfg *config.Config, userID uint, email string) (string, string, error) {
	expMinutes := cfg.JWT.ExpiresInMinutes
	if expMinutes == 0 {
		expMinutes = 15 // default 15 menit
	}

	refreshExpDays := cfg.JWT.RefreshExpiresInDays
	if refreshExpDays == 0 {
		refreshExpDays = 7 // default 7 hari
	}

	// Access Token
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "hrm-app",
			Subject:   "access_token",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, refreshExpDays)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "hrm-app",
			Subject:   "refresh_token",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(cfg *config.Config, tokenStr string) (*Claims, error) {
	return validateToken(cfg, tokenStr, "access_token")
}

func ValidateRefreshToken(cfg *config.Config, tokenStr string) (*Claims, error) {
	return validateToken(cfg, tokenStr, "refresh_token")
}

func validateToken(cfg *config.Config, tokenStr string, expectedSubject string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if expectedSubject != "" && claims.Subject != expectedSubject {
			return nil, errors.New("invalid token subject")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
