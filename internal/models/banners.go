package models

import "time"

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

type Banners struct {
	BannerId  int64     `json:"banner_id"`
	FeatureId int64     `json:"feature_id"`
	TagIds    []int64   `json:"tag_ids"`
	Content   Content   `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserBannersRequest struct {
	TagId          int64 `json:"tag_id"`
	FeatureId      int64 `json:"feature_id"`
	UseLastVersion bool  `json:"use_last_version,omitempty"`
}

type GetAllBannersRequest struct {
	FeatureId int64 `json:"feature_id,omitempty"`
	TagId     int64 `json:"tag_id,omitempty"`
	Limit     int64 `json:"limit,omitempty"`
	Offset    int64 `json:"offset,omitempty"`
}

type GetAllBannersResponse struct {
	BannerId  int64     `json:"banner_id"`
	FeatureId int64     `json:"feature_id"`
	TagIds    []int64   `json:"tag_ids"`
	Content   Content   `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBannerRequest struct {
	FeatureId int64   `json:"feature_id"`
	TagIds    []int64 `json:"tag_ids"`
	Content   Content `json:"content"`
	IsActive  bool    `json:"is_active"`
}

type UpdateBannerRequest struct {
	FeatureId int64   `json:"feature_id,omitempty"`
	TagIds    []int64 `json:"tag_ids,omitempty"`
	Content   Content `json:"content,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}
