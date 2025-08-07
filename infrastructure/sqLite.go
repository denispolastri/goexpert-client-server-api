package infrastructure

import (
	"goexpert-client-server-api/types"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewSqliteDb() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("infrastructure/db/sqlite.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&types.Dollar{})
	return db, nil
}
