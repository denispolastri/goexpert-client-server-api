package infrastructure

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewSqliteDb() {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate()
}
