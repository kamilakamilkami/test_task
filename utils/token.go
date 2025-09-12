package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"project/dto"
)

var (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

func GenerateTokens(email, role, secretKey string) (dto.TokenResponse, error) {
	accessToken, accessExpiry, err := generateAccessToken(email, role, secretKey)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	refreshToken, refreshExpiry, err := generateRefreshToken(email, role, secretKey)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}, nil
}

func generateAccessToken(email, role, secretKey string) (string, time.Time, error) {
	expiry := time.Now().Add(accessTokenTTL)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = expiry.Unix()

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiry, nil
}

func generateRefreshToken(email, role, secretKey string) (string, time.Time, error) {
	expiry := time.Now().Add(refreshTokenTTL)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["role"] = role
	claims["type"] = "refresh"
	claims["exp"] = expiry.Unix()

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiry, nil
}

func ParseToken(tokenString, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
func GenerateAccessTokenWithExpiry(email, role, secretKey string) (string, time.Time, error) {
	expiry := time.Now().Add(accessTokenTTL)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = expiry.Unix()

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expiry, nil
}
