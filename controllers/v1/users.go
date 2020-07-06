package v1

import (
	"net/http"

	"github.com/suarezgary/GolangApi/models"
	"github.com/suarezgary/GolangApi/utils/jsonhttp"
	"github.com/suarezgary/GolangApi/utils/queryparams"
)

// GetUsers GetUsers
func GetUsers(w http.ResponseWriter, r *http.Request) {
	limit, offset := queryparams.GetLimitOffsetQueryParametersDefaults(r)

	users, err := models.GetUsers(limit, offset)
	if err != nil {
		jsonhttp.JSONNotFoundError(w, "Error fetching gophers", "")
		return
	}

	jsonhttp.JSONSuccess(w, users, "Successfully queried gophers")
}

// CreateUser CreateUser
func CreateUser(w http.ResponseWriter, r *http.Request) {
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
