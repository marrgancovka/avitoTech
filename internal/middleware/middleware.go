package middleware

import (
	"avitoTech/internal/utils/jwter"
	"avitoTech/internal/utils/responser"
	"context"
	"errors"
	"net/http"
	"time"
)

// CookieName represents the name of the JWT cookie.
const CookieName = "jwt-banner-service"

// JwtMiddleware is a middleware function that handles JWT authentication.
func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil {
			responser.WriteError(w, http.StatusUnauthorized, errors.New("cookie not found"))
			return
		}
		token := cookie.Value
		claims, err := jwter.ParseToken(token)
		if err != nil {
			responser.WriteError(w, http.StatusUnauthorized, errors.New("invalid token: "+err.Error()))
			return
		}

		timeExp, err := claims.Claims.GetExpirationTime()
		if err != nil {
			responser.WriteError(w, http.StatusUnauthorized, errors.New("invalid token: "+err.Error()))
			return
		}
		if timeExp.Before(time.Now()) {
			responser.WriteError(w, http.StatusUnauthorized, errors.New("token is expired"))
			return
		}

		_, isAdmin, err := jwter.ParseClaims(claims)
		if err != nil {
			responser.WriteError(w, http.StatusUnauthorized, errors.New("error in parse token: "+err.Error()))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), CookieName, isAdmin))

		next.ServeHTTP(w, r)
	})
}
