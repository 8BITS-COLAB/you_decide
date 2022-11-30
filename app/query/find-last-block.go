package query

import (
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"gorm.io/gorm"
)

type FindLastBlockQuery struct {
	db *gorm.DB
}

func NewFindLastBlockQuery(db *gorm.DB) *FindLastBlockQuery {
	return &FindLastBlockQuery{
		db: db,
	}
}

func (q *FindLastBlockQuery) Exec() (*entity.Block, error) {
	var block entity.Block

	if err := q.db.Last(&block).Error; err != nil {
		return nil, err
	}

	return &block, nil
}
