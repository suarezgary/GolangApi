package models

import (
	"encoding/base64"
	"errors"
	"strings"

	"crypto/sha1"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

//User - User Model
type User struct {
	gorm.Model
	FullName       string         `json:"fullName"`
	Deleted        bool           `json:"-"`
	HashedPassword string         `json:"-" gorm:"type:varchar(500);"`
	Salt           string         `json:"-"`
	Address        string         `json:"address"`
	Email          string         `json:"email"`
	Password       string         `json:"password" gorm:"-"`
	NewPassword    string         `json:"newPassword" gorm:"-"`
	Notifications  []Notification `json:"notifications"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	uuidString := uuid.New().String()
	finalSalt := uuidString + user.Email + user.Password
	hasher := sha1.New()
	hasher.Write([]byte(finalSalt))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	scope.SetColumn("Salt", uuidString)
	scope.SetColumn("Deleted", false)
	return scope.SetColumn("HashedPassword", sha)
}

// FindByID - Find User by Id
func (user *User) FindByID() error {
	return db.Where("id = ?", user.ID).First(&user).Error
}

//Create Create User
func (user *User) Create() error {
	return db.Create(&user).Error
}

//UpdateMeta - Update User Object
func (user *User) UpdateMeta() error {
	return db.Table("user").Where("id = ?", user.ID).Updates(map[string]interface{}{
		"full_name": user.FullName,
		"address":   user.Address,
	}).Error
}

//GetUsers - Get Gophers
func GetUsers(limit, offset int64) ([]User, error) {
	var users []User
	if err := db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// FindUserByEmail - Find User by Email
func FindUserByEmail(email string) User {
	var user User
	if db.Where("email = ?", email).First(&user).RecordNotFound() {
		return user
	}
	return user
}

//ValidateUserModel ValidateUserModel
func (user *User) ValidateUserModel() error {
	if len(strings.TrimSpace(user.Password)) == 0 {
		return errors.New("Password can't be empty")
	}

	if len(strings.TrimSpace(user.Email)) == 0 {
		return errors.New("Email can't be empty")
	}

	if len(strings.TrimSpace(user.FullName)) == 0 {
		return errors.New("Name can't be empty")
	}
	return nil
}

//ValidateNewPassword ValidateNewPassword
func (user *User) ValidateNewPassword() error {
	if len(strings.TrimSpace(user.Password)) == 0 {
		return errors.New("Password can't be empty")
	}

	if len(strings.TrimSpace(user.Email)) == 0 {
		return errors.New("Email can't be empty")
	}

	if len(strings.TrimSpace(user.NewPassword)) == 0 {
		return errors.New("New Password can't be empty")
	}
	return nil
}

//ValidatePass ValidatePass
func (user *User) ValidatePass() bool {
	finalSalt := user.Salt + user.Email + user.Password
	hasher := sha1.New()
	hasher.Write([]byte(finalSalt))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha == user.HashedPassword
}

//ChangePass ChangePass
func (user *User) ChangePass(newPassword string) error {
	user.Password = newPassword
	user.Salt = uuid.New().String()
	finalSalt := user.Salt + user.Email + user.Password
	hasher := sha1.New()
	hasher.Write([]byte(finalSalt))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	user.HashedPassword = sha

	return db.Table("user").Where("id = ?", user.ID).Updates(map[string]interface{}{
		"hashed_password": user.HashedPassword,
		"salt":            user.Salt,
	}).Error
}
