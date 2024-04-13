package banners

type BannerUsecase interface {
}

type BannerRepo interface {
	GetUserBanner() error
	GetBanner() error
	CreateBanner() error
	UpdateBanner() error
	DeleteBanner() error
}
