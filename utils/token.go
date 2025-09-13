package utils

import (
	"project/dto"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(email, role, userId, secret string) (dto.TokenResponse, error) {
	accessToken, accessExpiry, err := GenerateAccessTokenWithExpiry(email, role, userId, secret)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	refreshToken, _ := GenerateRefreshToken(email, role, userId, secret)

	return dto.TokenResponse{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: time.Now().Add(7 * 24 * time.Hour),
	}, nil
}

func GenerateAccessTokenWithExpiry(email, role, userId, secret string) (string, time.Time, error) {
	exp := time.Now().Add(15 * time.Minute)
	claims := jwt.MapClaims{
		"id": userId,
		"email": email,
		"role":  role,
		"exp":   exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	return t, exp, err
}

func GenerateRefreshToken(email, role, userId, secret string) (string, error) {
	exp := time.Now().Add(7 * 24 * time.Hour)
	claims := jwt.MapClaims{
		"id": userId,
		"email": email,
		"role":  role,
		"exp":   exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
