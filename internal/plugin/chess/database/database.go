package database

import (
	"github.com/aimerneige/yukichan-bot/internal/plugin/chess/database/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a global variable for database
var DB *gorm.DB

// InitDatabase init database
func InitDatabase(filePath string) {
	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&model.ELO{}, &model.PGN{})
	if err != nil {
		panic(err)
	}

	DB = db
}

// GetDB get database
func GetDB() *gorm.DB {
	return DB
}
