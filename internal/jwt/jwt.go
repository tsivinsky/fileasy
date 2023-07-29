package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	AccessTokenLifeTime  = time.Hour * 24      // 24 hours
	RefreshTokenLifeTime = time.Hour * 24 * 30 // 30 days
)

var (
	JWTSecret = os.Getenv("JWT_SECRET")
)

type AuthTokenClaims struct {
	jwt.RegisteredClaims

	UserId uint `json:"userId"`
}

func GetKey(token *jwt.Token) (interface{}, error) {
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, errors.New("Invalid token")
	}

	return []byte(JWTSecret), nil
}

func GenerateAccessToken(userId uint) (string, error) {
	claims := AuthTokenClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenLifeTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(JWTSecret))

	return ss, err
}

func GenerateRefreshToken(userId uint) (string, error) {
	claims := AuthTokenClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenLifeTime)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte(JWTSecret))

	return ss, err
}

func GenerateBothTokens(userId uint) (string, string, error) {
	accessToken, err := GenerateAccessToken(userId)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(userId)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateRefreshToken(refreshToken string) (uint, error) {
	var claims AuthTokenClaims

	token, err := jwt.ParseWithClaims(refreshToken, &claims, GetKey)
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("Invalid refresh token")
	}

	return claims.UserId, nil
}

func ValidateAccessToken(accessToken string) (uint, error) {
	var claims AuthTokenClaims

	token, err := jwt.ParseWithClaims(accessToken, &claims, GetKey)
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("Invalid access token")
	}

	return claims.UserId, nil
}
