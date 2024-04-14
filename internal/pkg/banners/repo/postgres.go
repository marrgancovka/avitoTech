package repo

import (
	"avitoTech/internal/models"
	"database/sql"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) GetUserBanner(request *models.GetUserBannersRequest) (*models.Content, error) {

}
func (r *PostgresRepo) GetBanners(*models.GetAllBannersRequest) (*models.GetAllBannersResponse, error) {

}
func (r *PostgresRepo) CreateBanner(*models.CreateBannerRequest) (int64, error) {

}
func (r *PostgresRepo) UpdateBanner(int64, *models.CreateBannerRequest) error {

}
func (r *PostgresRepo) DeleteBanner(int64) error {

}
