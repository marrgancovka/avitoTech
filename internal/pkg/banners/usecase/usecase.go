package usecase

import (
	"avitoTech/internal/models"
	"avitoTech/internal/pkg/banners"
	"context"
	"encoding/json"
	"strconv"
)

type UseCase struct {
	r     banners.BannerRepo
	cache banners.BannerCache
}

func NewUseCase(r banners.BannerRepo, cache banners.BannerCache) *UseCase {
	return &UseCase{r: r, cache: cache}
}

func (uc *UseCase) GetUserBanner(ctx context.Context, request *models.GetUserBannersRequest) (*models.Content, error) {
	contentStruct := &models.Content{}
	key := strconv.FormatInt(request.TagId, 10) + strconv.FormatInt(request.FeatureId, 10)
	if !request.UseLastVersion {
		dataByte, ok := uc.cache.GetBanner(ctx, key)
		if ok {
			err := json.Unmarshal(dataByte, &contentStruct)
			if err != nil {
				return nil, err
			}
			return contentStruct, nil
		}
	}
	content, err := uc.r.GetUserBanner(ctx, request.TagId, request.FeatureId)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(content, contentStruct); err != nil {
		return nil, err
	}

	if err := uc.cache.SetBanner(ctx, key, content); err != nil {
		return nil, err
	}

	return contentStruct, nil
}
func (uc *UseCase) GetBanners(ctx context.Context, request *models.GetAllBannersRequest) ([]*models.GetAllBannersResponse, error) {
	bannersList, err := uc.r.GetBanners(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, banner := range bannersList {
		if !banner.IsActive {
			continue
		}
		content, err := json.Marshal(banner.Content)
		if err != nil {
			continue
		}
		for _, tagId := range banner.TagIds {
			key := strconv.FormatInt(tagId, 10) + strconv.FormatInt(banner.FeatureId, 10)
			if err = uc.cache.SetBanner(ctx, key, content); err != nil {
				continue
			}
		}

	}

	return bannersList, nil
}
func (uc *UseCase) CreateBanner(ctx context.Context, request *models.CreateBannerRequest) (int64, error) {
	bannerId, err := uc.r.CreateBanner(ctx, request)
	if err != nil {
		return 0, err
	}
	return bannerId, nil
}
func (uc *UseCase) UpdateBanner(ctx context.Context, id int64, request *models.UpdateBannerRequest) error {

	err := uc.r.UpdateBanner(ctx, id, request)
	return err
}
func (uc *UseCase) DeleteBanner(ctx context.Context, id int64) error {
	err := uc.r.DeleteBanner(ctx, id)
	return err
}
