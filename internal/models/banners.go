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
	TagId          int64  `json:"tag_id"`
	FeatureId      int64  `json:"feature_id"`
	UseLastVersion bool   `json:"use_last_version"`
	Token          string `json:"token"`
}

type GetAllBannersRequest struct {
	Token     string `json:"token"`
	FeatureId int64  `json:"feature_id"`
	TagId     int64  `json:"tag_id"`
	Limit     int64  `json:"limit"`
	Offset    int64  `json:"offset"`
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
