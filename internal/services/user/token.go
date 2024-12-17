package user

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const SecretKey = "jiaofupi"

func GenerateToken(ctx context.Context, userId string, t time.Duration) (token string, err error) {
	if t == 0 {
		t = 48 * time.Hour
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(t).Unix(),
	})

	token, err = accessToken.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func IsTokenExpiringSoon(ctx context.Context, token string) bool {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return false
	}
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Add(5*time.Minute).Unix() > int64(exp) {
				return true
			}
		}
	}
	return false
}

func ValidateToken(ctx context.Context, token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", ErrInvalidToken.WithError(err)
	}

	var userId string
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return "", ErrInvalidToken.WithError(fmt.Errorf("token has expired"))
			}
		} else {
			return "", ErrInvalidToken.WithError(fmt.Errorf("exp claim is not present or invalid"))
		}
		userId = claims["userId"].(string)
	} else {
		return "", ErrInvalidToken.WithError(fmt.Errorf("Invalid claims"))
	}
	return userId, nil
}
