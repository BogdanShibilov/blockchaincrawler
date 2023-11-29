package repository

import (
	"math"

	"gorm.io/gorm"
)

func paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (b *Block) GetTotalPagesFor(entity any, pageSize int) (totalPages int32) {
	var totalRows int64
	b.main.DB.Model(entity).Count(&totalRows)

	return int32(math.Ceil(float64(totalRows) / float64(pageSize)))
}
