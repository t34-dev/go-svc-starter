package auth_service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

func generateTokens(userID uuid.UUID, sessionID uuid.UUID) (string, string, error) {
	// Access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID.String(),
		"session_id": sessionID.String(),
		"exp":        time.Now().Add(15 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID.String(),
		"session_id": sessionID.String(),
		"exp":        time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 дней
	})

	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, nil
}

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
