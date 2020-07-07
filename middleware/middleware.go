package middleware

import (
	"net/http"
	"strings"

	"github.com/suarezgary/GolangApi/config"
	"github.com/suarezgary/GolangApi/models"
	"github.com/suarezgary/GolangApi/reqctx"
	"github.com/suarezgary/GolangApi/utils/htmlhttp"
	"github.com/suarezgary/GolangApi/utils/jsonhttp"
	"github.com/suarezgary/GolangApi/utils/jwtutil"
)

//Log - Logger
var Log = config.Cfg().GetLogger()

//GetContext - Get Context
func GetContext(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var passItOn = func() {
			handler.ServeHTTP(w, r)
		}
		var handleNonEmptyToken = func(value string) {
			var user models.User
			userToken, err := jwtutil.ExtractTokenMetadata(value)
			if err != nil {
				passItOn()
				return
			}
			user.ID = userToken.ID
			user.FullName = userToken.FullName
			// Update our request
			r = r.WithContext(reqctx.AddCurrentUserToContext(r, user))
			// Move on to the next middleware
			passItOn()
		}

		if c, err := r.Cookie(config.Cfg().TokenCookieName); err == nil {
			if c.Value != "" {
				handleNonEmptyToken(c.Value)
			} else {
				passItOn()
			}
		} else if ah := r.Header.Get("Authorization"); ah != "" {
			if len(ah) > 6 && strings.ToUpper(ah[0:7]) == "BEARER " {
				val := ah[7:]
				if val != "" {
					handleNonEmptyToken(val)
				} else {
					passItOn()
				}
			} else {
				passItOn()
			}
		} else {
			passItOn()
		}
	}
}

// RequireAPIKey RequireAPIKey
func RequireAPIKey(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var respondUnauthorized = func() {
			jsonhttp.JSONDetailed(w, jsonhttp.APIResponse{Message: "Unauthorized", Debug: "Invalid or missing API Key header/query parameter"}, http.StatusUnauthorized)
		}

		if r.URL.Query().Get("apiKey") != config.Cfg().ServiceAPIKey && r.Header.Get("x-api-key") != config.Cfg().ServiceAPIKey {
			respondUnauthorized()
			return
		}

		handler.ServeHTTP(w, r)
	}
}

// RequireUserForAPI RequireUserForAPI
func RequireUserForAPI(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var respondUnauthorized = func() {
			jsonhttp.JSONDetailed(w, jsonhttp.APIResponse{Message: "Unauthorized", Debug: "Invalid or missing access token header/cookie"}, http.StatusUnauthorized)
		}
		user, err := reqctx.GetCurrentUser(r)
		_ = user
		if err != nil {
			respondUnauthorized()
			return
		}

		handler.ServeHTTP(w, r)
	}
}

// RequireUserForView RequireUserForView
func RequireUserForView(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var respondUnauthorized = func() {
			htmlhttp.UnauthorizedErrorView(w, r)
		}
		user, err := reqctx.GetCurrentUser(r)

		_ = user
		if err != nil {
			respondUnauthorized()
			return
		}

		handler.ServeHTTP(w, r)
	}
}
