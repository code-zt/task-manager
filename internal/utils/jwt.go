package utils

import (
	"errors"
	"fmt"
	"task_manager/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	secretKey       = "secret_key"
	accessTokenExp  = 24 * time.Hour      // 24 часа для access токена
	refreshTokenExp = 30 * 24 * time.Hour // 30 дней для refresh токена
)

// CreateAccessToken создает access токен
func CreateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"name":  user.Username,
		"email": user.Email,
		"id":    user.ID.Hex(),
		"exp":   time.Now().Add(accessTokenExp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// CreateRefreshToken создает refresh токен
func CreateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":  user.ID.Hex(),
		"exp": time.Now().Add(refreshTokenExp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateAccessToken проверяет access токен
func ValidateAccessToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Извлекаем claims и создаём объект пользователя
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Проверяем срок действия токена
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return nil, errors.New("token has expired")
			}
		} else {
			return nil, errors.New("expiration claim not found")
		}

		// Создаём объект пользователя из claims
		user := &models.User{
			Username: claims["name"].(string),
			Email:    claims["email"].(string),
		}

		// Если ID хранится как string (Hex)
		if idStr, ok := claims["id"].(string); ok {
			if objectID, err := primitive.ObjectIDFromHex(idStr); err == nil {
				user.ID = objectID
			}
		}

		return user, nil
	}

	return nil, errors.New("invalid token claims")
}

// ValidateRefreshToken проверяет refresh токен
func ValidateRefreshToken(tokenString string) (primitive.ObjectID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return primitive.ObjectID{}, errors.New("invalid token")
	}

	// Извлекаем claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Проверяем срок действия токена
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return primitive.ObjectID{}, errors.New("token has expired")
			}
		} else {
			return primitive.ObjectID{}, errors.New("expiration claim not found")
		}

		// Возвращаем ID пользователя
		if idStr, ok := claims["id"].(string); ok {
			if objectID, err := primitive.ObjectIDFromHex(idStr); err == nil {
				return objectID, nil
			}
		}
	}

	return primitive.ObjectID{}, errors.New("invalid token claims")
}

// RefreshAccessToken обновляет access токен с использованием refresh токена
func RefreshAccessToken(refreshToken string) (string, error) {
	userID, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Здесь вы можете получить пользователя из базы данных по userID
	// user := getUser ByID(userID)

	// Для примера, создадим временного пользователя
	user := &models.User{
		ID: userID,
		// Заполните остальные поля, если необходимо
	}

	// Создаем новый access токен
	return CreateAccessToken(user)
}
