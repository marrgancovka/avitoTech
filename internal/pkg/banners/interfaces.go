package banners

import "avitoTech/internal/models"

type BannerUsecase interface {
	GetUserBanner(request *models.GetUserBannersRequest) (*models.Content, error)
	GetBanners(request *models.GetAllBannersRequest) (*models.GetAllBannersResponse, error)
	CreateBanner(request *models.CreateBannerRequest) (int64, error)
	UpdateBanner(int64, *models.CreateBannerRequest) error
	DeleteBanner(int64) error
}

type BannerRepo interface {
	GetUserBanner(request *models.GetUserBannersRequest) (*models.Content, error)
	GetBanners(*models.GetAllBannersRequest) (*models.GetAllBannersResponse, error)
	CreateBanner(*models.CreateBannerRequest) (int64, error)
	UpdateBanner(int64, *models.CreateBannerRequest) error
	DeleteBanner(int64) error
}

type BannerCache interface {
	GetBanner() error
	SetBanner() error
}
