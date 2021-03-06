package models

import (
	"github.com/jinzhu/gorm"
	"github.com/suarezgary/GolangApi/config"

	// Use the _ import syntax to ensure that the mysql init()
	// gets run and so that Go doesn't complain about an unused import
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Log - Logger
var Log = config.Cfg().GetLogger()
var db *gorm.DB

//Setup - Setup database
func Setup() (err error) {
	path := config.Cfg().DBPath
	db, err = gorm.Open("mysql", path)
	if err != nil {
		return err
	}

	// Some reasonable pool settings
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Log all our queries
	db.LogMode(true)

	db.SingularTable(true)

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Notification{})
	db.Model(&Notification{}).AddForeignKey("user_id", "user(id)", "CASCADE", "CASCADE")

	return nil
}
