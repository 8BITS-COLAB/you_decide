package query

import (
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"gorm.io/gorm"
)

type FindVoteQuery struct {
	db *gorm.DB
}

func NewFindVoteQuery(db *gorm.DB) *FindVoteQuery {
	return &FindVoteQuery{
		db: db,
	}
}

func (q *FindVoteQuery) Exec(where ...interface{}) (*entity.Vote, error) {
	var vote entity.Vote

	if err := q.db.First(&vote, where...).Error; err != nil {
		return nil, err
	}

	return &vote, nil
}
