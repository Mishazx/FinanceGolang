package repository

import (
	"gorm.io/gorm"
)

// BaseRepository to hold the common database instance
type BaseRepository struct {
	db *gorm.DB
}

// InitializeRepository initializes a new BaseRepository with a DB instance
func InitializeRepository(db *gorm.DB) BaseRepository {
	return BaseRepository{db: db}
}
