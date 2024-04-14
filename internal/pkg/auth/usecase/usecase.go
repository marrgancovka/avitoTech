package usecase

import (
	"avitoTech/internal/models"
	"avitoTech/internal/pkg/auth"
	"avitoTech/internal/utils/jwter"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/satori/uuid"
	"time"
)

type AuthUsecase struct {
	r auth.AuthRepo
}

func NewUsecase(r auth.AuthRepo) *AuthUsecase {
	return &AuthUsecase{r: r}
}

func (uc *AuthUsecase) SignUp(ctx context.Context, authData *models.AuthRequest) (string, time.Time, error) {
	newUser := &models.User{
		Id:           uuid.NewV4(),
		Username:     authData.Username,
		PasswordHash: hashString(authData.Password),
		IsAdmin:      false,
	}
	if err := uc.r.CreateUser(ctx, newUser); err != nil {
		return "", time.Time{}, err
	}
	token, exp, err := jwter.GenerateToken(newUser)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, exp, nil
}

func (uc *AuthUsecase) SignIn(ctx context.Context, authData *models.AuthRequest) (string, time.Time, error) {
	user, err := uc.r.GetUser(ctx, authData.Username)
	if err != nil {
		return "", time.Time{}, err
	}
	if hashString(authData.Password) != user.PasswordHash {
		return "", time.Time{}, errors.New("invalid password")
	}
	token, exp, err := jwter.GenerateToken(user)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, exp, nil
}

func hashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))

	hashedBytes := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}
