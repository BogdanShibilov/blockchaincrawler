package pagination

import (
	"math"

	"gorm.io/gorm"
)

func Paginate(value any, pgn *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	pgn.TotalRows = totalRows

	totalPages := int(math.Ceil(float64(totalRows) / float64(pgn.GetLimit())))
	pgn.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pgn.GetOffset()).Limit(pgn.GetLimit()).Order(pgn.GetSort())
	}
}
