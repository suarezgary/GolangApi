package v1

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/suarezgary/GolangApi/models"
	"github.com/suarezgary/GolangApi/utils/emailsender"
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
		return
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

	emailsender.SendWelcomeEmail(user.Email, user.FullName)

	user.Password = ""
	jsonhttp.JSONSuccess(w, user, "Successfully created User")
}

// Forgot Forgot
func Forgot(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := jsonhttp.DecodeJSONBody(w, r, &user)
	if err != nil {
		jsonhttp.JSONInternalError(w, "Error Forgot", err.Error())
		return
	}
	if len(strings.TrimSpace(user.Email)) == 0 {
		jsonhttp.JSONInternalError(w, "Email can't be empty", "")
		return
	}

	existingUser := models.FindUserByEmail(user.Email)
	if existingUser.ID == 0 {
		jsonhttp.JSONBadRequestError(w, "Email not valid", "")
		return
	}

	newPass := randStringBytesMaskImprSrc(8)
	err = existingUser.ChangePass(newPass)
	if err != nil {
		jsonhttp.JSONInternalError(w, "Error reseting password", err.Error())
		return
	}

	emailsender.SendForgotEmail(existingUser.Email, existingUser.FullName, existingUser.Email, newPass)

	user.Password = ""
	jsonhttp.JSONSuccess(w, nil, "Password Successfully Reset")
}

// ChangePass ChangePass
func ChangePass(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := jsonhttp.DecodeJSONBody(w, r, &user)
	if err != nil {
		jsonhttp.JSONInternalError(w, "Error Change Password", err.Error())
		return
	}
	err = user.ValidateNewPassword()
	if err != nil {
		jsonhttp.JSONInternalError(w, err.Error(), err.Error())
		return
	}

	existingUser := models.FindUserByEmail(user.Email)
	if existingUser.ID == 0 {
		jsonhttp.JSONBadRequestError(w, "Email not valid", "")
		return
	}

	err = existingUser.ChangePass(user.NewPassword)
	if err != nil {
		jsonhttp.JSONInternalError(w, "Error changing password", err.Error())
		return
	}

	emailsender.SendChangePass(existingUser.Email, existingUser.FullName)

	user.Password = ""
	jsonhttp.JSONSuccess(w, nil, "Password Successfully Reset")
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrc(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
