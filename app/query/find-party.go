package query

import (
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"gorm.io/gorm"
)

type FindPartyQuery struct {
	db *gorm.DB
}

func NewFindPartyQuery(db *gorm.DB) *FindPartyQuery {
	return &FindPartyQuery{
		db: db,
	}
}

func (q *FindPartyQuery) Exec(where ...interface{}) (*entity.Party, error) {
	var party entity.Party

	if err := q.db.First(&party, where...).Error; err != nil {
		return nil, err
	}

	return &party, nil
}
