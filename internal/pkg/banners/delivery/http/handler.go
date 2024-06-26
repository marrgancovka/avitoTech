package http

import (
	"avitoTech/internal/models"
	"avitoTech/internal/pkg/banners"
	"avitoTech/internal/utils/responser"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler struct {
	uc  banners.BannerUsecase
	log *logrus.Logger
}

func NewHandler(uc banners.BannerUsecase, log *logrus.Logger) *Handler {
	return &Handler{uc: uc, log: log}
}

func (h *Handler) GetUserBanners(w http.ResponseWriter, r *http.Request) {
	tagQuery := r.URL.Query().Get("tag_id")
	featureQuery := r.URL.Query().Get("feature_id")
	useLastRevisionQuery := r.URL.Query().Get("use_last_revision")
	useLastRevision := false

	tagId, err := strconv.ParseInt(tagQuery, 10, 64)
	if err != nil {
		h.log.Error("Error parsing tag_id")
		responser.WriteError(w, http.StatusBadRequest, errors.New("invalid tag_id"))
		return
	}
	featureId, err := strconv.ParseInt(featureQuery, 10, 64)
	if err != nil {
		h.log.Error("Error parsing feature_id")
		responser.WriteError(w, http.StatusBadRequest, errors.New("invalid feature_id"))
		return
	}
	useLastRevision, err = strconv.ParseBool(useLastRevisionQuery)
	if err != nil {
		h.log.Warn("Error parsing use_last_revision")
	}

	request := &models.GetUserBannersRequest{
		TagId:          tagId,
		FeatureId:      featureId,
		UseLastVersion: useLastRevision,
	}

	content, err := h.uc.GetUserBanner(r.Context(), request)
	if err != nil {
		h.log.Error("Error getting user banners: " + err.Error())
		responser.WriteError(w, http.StatusNotFound, err)
		return
	}
	responser.WriteJSON(w, http.StatusOK, content)
}

func (h *Handler) GetBanners(w http.ResponseWriter, r *http.Request) {
	tagQuery := r.URL.Query().Get("tag_id")
	featureQuery := r.URL.Query().Get("feature_id")
	limitQuery := r.URL.Query().Get("limit")
	offsetQuery := r.URL.Query().Get("offset")
	var tagId, featureId, limit, offset int64

	tagId, err := strconv.ParseInt(tagQuery, 10, 64)
	if err != nil {
		h.log.Warn("Error parsing tag_id")
	}
	featureId, err = strconv.ParseInt(featureQuery, 10, 64)
	if err != nil {
		h.log.Warn("Error parsing feature_id")
	}
	limit, err = strconv.ParseInt(limitQuery, 10, 64)
	if err != nil {
		h.log.Warn("Error parsing limit")
	}
	offset, err = strconv.ParseInt(offsetQuery, 10, 64)
	if err != nil {
		h.log.Warn("Error parsing offset")
	}
	request := &models.GetAllBannersRequest{
		FeatureId: featureId,
		TagId:     tagId,
		Limit:     limit,
		Offset:    offset,
	}

	bannersList, err := h.uc.GetBanners(r.Context(), request)
	if err != nil {
		h.log.Error("ErrorGetGetBanners", err)
		responser.WriteError(w, http.StatusNotFound, err)
		return
	}

	responser.WriteJSON(w, http.StatusOK, bannersList)
}

func (h *Handler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	data := &models.CreateBannerRequest{}
	if err := responser.ReadRequestData(r, &data); err != nil {
		h.log.Errorf("ErrorPostRequestBody: %v", err)
		responser.WriteError(w, http.StatusBadRequest, errors.New("Incorrect data format: "+err.Error()))
		return
	}
	bannerId, err := h.uc.CreateBanner(r.Context(), data)
	if err != nil {
		h.log.Error("Error in create banner: ", err.Error())
		responser.WriteError(w, http.StatusBadRequest, errors.New("Error in create banner: "+err.Error()))
		return
	}
	responser.WriteJSON(w, http.StatusOK, bannerId)
}

func (h *Handler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bannerIdQuery := vars["id"]

	bannerId, err := strconv.ParseInt(bannerIdQuery, 10, 64)
	if err != nil {
		h.log.Error("Error parsing banner id")
		responser.WriteError(w, http.StatusBadRequest, errors.New("invalid banner id"))
		return
	}

	data := &models.UpdateBannerRequest{}
	if err := responser.ReadRequestData(r, &data); err != nil {
		h.log.Errorf("ErrorPatchRequestBody: %v", err)
		responser.WriteError(w, http.StatusBadRequest, errors.New("Incorrect data format: "+err.Error()))
		return
	}

	if err := h.uc.UpdateBanner(r.Context(), bannerId, data); err != nil {
		h.log.Errorf("ErrorUpdateBanner: %v", err)
		responser.WriteError(w, http.StatusNotFound, err)
		return
	}
	responser.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bannerIdQuery := vars["id"]

	bannerId, err := strconv.ParseInt(bannerIdQuery, 10, 64)
	if err != nil {
		h.log.Error("Error parsing banner id")
		responser.WriteError(w, http.StatusBadRequest, errors.New("invalid banner id"))
		return
	}

	if err := h.uc.DeleteBanner(r.Context(), bannerId); err != nil {
		h.log.Errorf("ErrorDeleteBanner: %v", err)
		responser.WriteError(w, http.StatusNotFound, err)
		return
	}
	responser.WriteJSON(w, http.StatusNoContent, nil)
}
