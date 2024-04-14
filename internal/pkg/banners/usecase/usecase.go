package usecase

import (
	"avitoTech/internal/models"
	"avitoTech/internal/pkg/banners"
)

type UseCase struct {
	r     banners.BannerRepo
	cache banners.BannerCache
}

func NewUseCase(r banners.BannerRepo, cache banners.BannerCache) *UseCase {
	return &UseCase{r: r, cache: cache}
}

func (uc *UseCase) GetUserBanner(request *models.GetUserBannersRequest) (*models.Content, error) {

}
func (uc *UseCase) GetBanners(request *models.GetAllBannersRequest) (*models.GetAllBannersResponse, error) {

}
func (uc *UseCase) CreateBanner(request *models.CreateBannerRequest) (int64, error) {

}
func (uc *UseCase) UpdateBanner(int64, *models.CreateBannerRequest) error {

}
func (uc *UseCase) DeleteBanner(int64) error {

}
