package jwter

import (
	"avitoTech/internal/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/satori/uuid"
	"net/http"
	"os"
	"time"
)

// GenerateToken returns a new JWT token for the given user.
func GenerateToken(user *models.User) (string, time.Time, error) {
	exp := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"is_admin": user.IsAdmin,
		"exp":      exp.Unix(),
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", time.Now(), err
	}
	return tokenStr, exp, nil
}

// ParseToken parses the provided JWT token string and returns the parsed token.
func ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

// ParseClaims parses the user ID from the JWT token claims.
func ParseClaims(claims *jwt.Token) (uuid.UUID, bool, error) {
	payloadMap, ok := claims.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, false, errors.New("invalid claims")
	}
	idStr, ok := payloadMap["id"].(string)
	if !ok {
		return uuid.Nil, false, errors.New("incorrect id")
	}
	id, err := uuid.FromString(idStr)
	if err != nil {
		return uuid.Nil, false, errors.New("incorrect id")
	}
	isAdmin, ok := payloadMap["is_admin"].(bool)
	if !ok {
		return uuid.Nil, false, errors.New("incorrect level")
	}

	return id, isAdmin, nil
}

// TokenCookie creates a new cookie for storing the authentication token.
func TokenCookie(name, token string, exp time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    token,
		Expires:  exp,
		Path:     "/",
		HttpOnly: true,
	}
}
