package database

import (
	"gofit-api/models"

	"gorm.io/gorm"
)

func PaginatedQuery(page *models.Pages) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(page.Offset).Limit(page.PageSize)
	}
}
