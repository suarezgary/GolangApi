package models

import "github.com/jinzhu/gorm"

const (
	email = iota // 0
	push  = iota // 1
	both  = iota // 2
)

//Notification Notification
type Notification struct {
	gorm.Model
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Type     uint32 `json:"type"`
	SendType uint32 `json:"sendType"`
	Sended   bool   `json:"sended"`
	Read     bool   `json:"read"`
	InfoJSON string `json:"infoJSON"`
	UserID   uint   `json:"userId"`
}

// FindByID - Find User by Id
func (notification *Notification) FindByID() error {
	return db.Where("id = ?", notification.ID).First(&notification).Error
}

//Create Create User
func (notification *Notification) Create() error {
	return db.Create(&notification).Error
}

//UpdateMeta - Update User Object
func (notification *Notification) UpdateMeta() error {
	return db.Table("notification").Where("id = ?", notification.ID).Updates(map[string]interface{}{
		"sended": notification.Sended,
		"read":   notification.Read,
	}).Error
}

//GetNotifications - Get Notifications
func GetNotifications(limit, offset int64, userID uint64) ([]Notification, error) {
	var notification []Notification
	if err := db.Limit(limit).Offset(offset).Where("user_id = ?", userID).Find(&notification).Error; err != nil {
		return nil, err
	}
	return notification, nil
}
