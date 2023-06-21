package database

import (
	"gofit-api/configs"
	"gofit-api/models"

	"gorm.io/gorm"
)

func PaginatedQuery(page *models.Pages) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(page.Offset).Limit(page.PageSize)
	}
}

func CountTotalData(tableName string) int {
	var totalData int64
	configs.DB.Table(tableName).Where("deleted_at IS NULL").Count(&totalData)
	return int(totalData)
}
