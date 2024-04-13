package auth

import (
	"avitoTech/internal/models"
	"context"
	"github.com/satori/uuid"
	"time"
)

type AuthUsecase interface {
	SignUp(context.Context, *models.AuthRequest) (uuid.UUID, string, time.Time, error)
	SignIn(context.Context, *models.AuthRequest) (*models.User, string, time.Time, error)
}

type AuthRepo interface {
	CreateUser(ctx context.Context, newUser *models.User) error
	GetUser(ctx context.Context, username string) (*models.User, error)
}
