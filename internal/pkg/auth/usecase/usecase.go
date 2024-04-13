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

func (uc *AuthUsecase) SignUp(ctx context.Context, authData *models.AuthRequest) (uuid.UUID, string, time.Time, error) {
	newUser := &models.User{
		Id:           uuid.NewV4(),
		Username:     authData.Username,
		PasswordHash: hashString(authData.Password),
		IsAdmin:      false,
	}
	if err := uc.r.CreateUser(ctx, newUser); err != nil {
		return uuid.Nil, "", time.Time{}, err
	}
	token, exp, err := jwter.GenerateToken(newUser)
	if err != nil {
		return uuid.Nil, "", time.Time{}, err
	}
	return newUser.Id, token, exp, nil
}

func (uc *AuthUsecase) SignIn(ctx context.Context, authData *models.AuthRequest) (*models.User, string, time.Time, error) {
	user, err := uc.r.GetUser(ctx, authData.Username)
	if err != nil {
		return nil, "", time.Time{}, err
	}
	if authData.Password != user.PasswordHash {
		return nil, "", time.Time{}, errors.New("invalid password")
	}
	token, exp, err := jwter.GenerateToken(user)
	if err != nil {
		return nil, "", time.Time{}, err
	}
	return user, token, exp, nil
}

func hashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))

	hashedBytes := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString
}
