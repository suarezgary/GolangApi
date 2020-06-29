package models

import (
	"time"
)

//Gopher - Gopher Model
type Gopher struct {
	ID        uint64    `json:"id"`
	FullName  string    `json:"fullName"`
	Headline  string    `json:"headline"`
	AvatarURL string    `json:"avatarUrl"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

const getGophersPaginatedQuery = `SELECT * FROM gopher ORDER BY created_at DESC LIMIT ? OFFSET ?;`

// FindByID - Find Gopher by Id
func (gopher *Gopher) FindByID() error {
	return db.Where("id = ?", gopher.ID).First(&gopher).Error
}

//Create Create Gopher
func (gopher *Gopher) Create() error {
	return db.Create(&gopher).Error
}

//UpdateMeta - Update Gopher Object
func (gopher *Gopher) UpdateMeta() error {
	return db.Table("gopher").Where("id = ?", gopher.ID).Updates(map[string]interface{}{
		"full_name":  gopher.FullName,
		"headline":   gopher.Headline,
		"avatar_url": gopher.AvatarURL,
	}).Error
}

//GetGophers - Get Gophers
func GetGophers(limit, offset int64) ([]Gopher, error) {
	var gophers []Gopher
	// Technically this query could go through GORM natively, but just showing off raw SQL query functionality!
	if err := db.Raw(getGophersPaginatedQuery, limit, offset).Scan(&gophers).Error; err != nil {
		return nil, err
	}
	return gophers, nil
}
