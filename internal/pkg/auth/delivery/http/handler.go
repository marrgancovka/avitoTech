package http

import (
	"avitoTech/internal/middleware"
	"avitoTech/internal/models"
	"avitoTech/internal/pkg/auth"
	"avitoTech/internal/utils/jwter"
	"avitoTech/internal/utils/responser"
	"errors"
	"net/http"
)

type AuthHandler struct {
	uc auth.AuthUsecase
}

func NewHandler(uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	data := models.AuthRequest{}

	if err := responser.ReadRequestData(r, &data); err != nil {
		responser.WriteError(w, http.StatusBadRequest, errors.New("Incorrect data format: "+err.Error()))
		return
	}

	token, exp, err := h.uc.SignIn(r.Context(), &data)
	if err != nil {
		responser.WriteError(w, http.StatusBadRequest, errors.New("Incorrect data format: "+err.Error()))
		return
	}
	http.SetCookie(w, jwter.TokenCookie(middleware.CookieName, token, exp))
	responser.WriteJSON(w, http.StatusCreated, token)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := models.AuthRequest{}

	if err := responser.ReadRequestData(r, &data); err != nil {
		responser.WriteError(w, http.StatusBadRequest, errors.New("Incorrect data format: "+err.Error()))
		return
	}

	token, exp, err := h.uc.SignUp(r.Context(), &data)
	if err != nil {
		responser.WriteError(w, http.StatusBadRequest, errors.New("Incorrect data format: "+err.Error()))
		return
	}
	http.SetCookie(w, jwter.TokenCookie(middleware.CookieName, token, exp))
	responser.WriteJSON(w, http.StatusCreated, token)
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  middleware.CookieName,
		Value: "",
		Path:  "/",
	})
	responser.WriteJSON(w, http.StatusOK, responser.MessageResponse{Message: "success logout"})
}
