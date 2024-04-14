package repo

import (
	"avitoTech/internal/models"
	"context"
	"database/sql"
)

type AuthRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(ctx context.Context, user *models.User) error {
	insert := `INSERT INTO "user" (id, username, password_hash) VALUES ($1, $2, $3)`
	if _, err := r.db.ExecContext(ctx, insert, user.Id, user.Username, user.PasswordHash); err != nil {
		return err
	}
	return nil
}

func (r *AuthRepo) GetUser(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	getUser := `SELECT id, username, password_hash, is_admin FROM "user" WHERE username=$1`
	if err := r.db.QueryRowContext(ctx, getUser, username).Scan(&user.Id, &user.Username, &user.PasswordHash, &user.IsAdmin); err != nil {
		return nil, err
	}
	return user, nil
}
