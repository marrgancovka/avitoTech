package banners

import (
	"avitoTech/internal/models"
	"context"
)

type BannerUsecase interface {
	GetUserBanner(context.Context, *models.GetUserBannersRequest) (*models.Content, error)
	GetBanners(context.Context, *models.GetAllBannersRequest) ([]*models.GetAllBannersResponse, error)
	CreateBanner(context.Context, *models.CreateBannerRequest) (int64, error)
	UpdateBanner(context.Context, int64, *models.UpdateBannerRequest) error
	DeleteBanner(context.Context, int64) error
}

type BannerRepo interface {
	GetUserBanner(context.Context, int64, int64) ([]byte, error)
	GetBanners(context.Context, *models.GetAllBannersRequest) ([]*models.GetAllBannersResponse, error)
	CreateBanner(context.Context, *models.CreateBannerRequest) (int64, error)
	UpdateBanner(context.Context, int64, *models.UpdateBannerRequest) error
	DeleteBanner(context.Context, int64) error
}

type BannerCache interface {
	GetBanner(context.Context, string) ([]byte, bool)
	SetBanner(context.Context, string, []byte) error
}
