package reqctx

import (
	"context"
	"errors"
	"net/http"

	"github.com/suarezgary/GolangApi/models"
	"github.com/suarezgary/GolangApi/utils/htmlhttp"
	"github.com/suarezgary/GolangApi/utils/jsonhttp"
)

const currentUserContextKey = "req_usr"

var unauthorizedError = jsonhttp.APIResponse{Message: "Unauthorized", Success: false, Debug: "Token/account error"}

//AddCurrentUserToContext - Add Current User to Context
func AddCurrentUserToContext(r *http.Request, user models.User) context.Context {
	return context.WithValue(r.Context(), currentUserContextKey, user)
}

//GetCurrentUser - Get Current User from Context
func GetCurrentUser(r *http.Request) (models.User, error) {
	user := r.Context().Value(currentUserContextKey)
	switch user.(type) {
	case nil:
		return models.User{}, errors.New("current user not found")
	default:
		return user.(models.User), nil
	}
}

//GetCurrentUserAndCatchForAPI - Get Current USer And Catch For Api
func GetCurrentUserAndCatchForAPI(w http.ResponseWriter, r *http.Request) (models.User, error) {
	var user models.User
	user, err := GetCurrentUser(r)
	if err != nil {
		jsonhttp.JSONWriter(w, unauthorizedError, http.StatusUnauthorized)
		return user, err
	}
	return user, nil
}

//GetCurrentUserAndCatchForView - Get current USer and Catch for View
func GetCurrentUserAndCatchForView(w http.ResponseWriter, r *http.Request) (models.User, error) {
	var user models.User
	user, err := GetCurrentUser(r)
	if err != nil {
		htmlhttp.UnauthorizedErrorView(w, r)
		return user, err
	}
	return user, nil
}
