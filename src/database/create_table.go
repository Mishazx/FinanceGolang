package database

import (
	"FinanceGolang/src/model"

	"gorm.io/gorm"
)

func CreateTables(db *gorm.DB) error {
	return db.AutoMigrate(&model.User{}, &model.Account{}, &model.Card{})
}
