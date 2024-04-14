package middleware

import (
	"avitoTech/internal/utils/jwter"
	"avitoTech/internal/utils/responser"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// CookieName represents the name of the JWT cookie.
const CookieName = "jwt-banner-service"

// Auth is a middleware function that handles JWT authentication.
func Auth(next http.Handler, log *logrus.Logger, forAdmin bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		claims, err := jwter.ParseToken(token)
		if err != nil {
			log.Error("Error parsing token: ", err.Error())
			responser.WriteError(w, http.StatusUnauthorized, errors.New("invalid token: "+err.Error()))
			return
		}

		timeExp, err := claims.Claims.GetExpirationTime()
		if err != nil {
			log.Error("invalid token: ", err.Error())
			responser.WriteError(w, http.StatusUnauthorized, errors.New("invalid token: "+err.Error()))
			return
		}
		if timeExp.Before(time.Now()) {
			log.Error("token is expired")
			responser.WriteError(w, http.StatusUnauthorized, errors.New("token is expired"))
			return
		}

		_, isAdmin, err := jwter.ParseClaims(claims)
		if err != nil {
			log.Error("Error parsing claims: ", err.Error())
			responser.WriteError(w, http.StatusUnauthorized, errors.New("error in parse token: "+err.Error()))
			return
		}

		if !isAdmin && forAdmin {
			responser.WriteError(w, http.StatusForbidden, errors.New("no access rights"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
