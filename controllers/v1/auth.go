package v1

import (
	"net/http"
	"strings"

	"github.com/suarezgary/GolangApi/models"
	"github.com/suarezgary/GolangApi/utils/jsonhttp"
	"github.com/suarezgary/GolangApi/utils/jwtutil"
)

// Login Login
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := jsonhttp.DecodeJSONBody(w, r, &user)
	if err != nil {
		jsonhttp.JSONInternalError(w, "Error on Request", err.Error())
		return
	}
	if len(strings.TrimSpace(user.Email)) == 0 || len(strings.TrimSpace(user.Password)) == 0 {
		jsonhttp.JSONBadRequestError(w, "Email and Password required", "")
	}

	existingUser := models.FindUserByEmail(user.Email)
	if existingUser.ID == 0 {
		jsonhttp.JSONBadRequestError(w, "Email/Password not valid", "")
		return
	}

	existingUser.Password = user.Password
	if !existingUser.ValidatePass() {
		jsonhttp.JSONBadRequestError(w, "Email/Password not valid", "")
		return
	}

	token, err := jwtutil.CreateToken(existingUser)
	if err != nil {
		jsonhttp.JSONInternalError(w, "Login Error", err.Error())
		return
	}

	jsonhttp.JSONSuccess(w, token, "Login Success")
}

// SignUp SignUp
func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := jsonhttp.DecodeJSONBody(w, r, &user)
	if err != nil {
		jsonhttp.JSONInternalError(w, "Error Creating the User", err.Error())
		return
	}
	err = user.ValidateUserModel()
	if err != nil {
		jsonhttp.JSONInternalError(w, err.Error(), err.Error())
		return
	}

	existingUser := models.FindUserByEmail(user.Email)
	if existingUser.ID != 0 {
		jsonhttp.JSONBadRequestError(w, "Email not available", "")
		return
	}

	err = user.Create()
	if err != nil {
		jsonhttp.JSONInternalError(w, "Error Creating the User", err.Error())
		return
	}
	user.Password = ""
	jsonhttp.JSONSuccess(w, user, "Successfully created User")
}
