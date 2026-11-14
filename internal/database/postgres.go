package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect membuka koneksi ke Postgres lewat GORM
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}