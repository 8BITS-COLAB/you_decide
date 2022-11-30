package query

import (
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"gorm.io/gorm"
)

type FindVoterQuery struct {
	db *gorm.DB
}

func NewFindVoterQuery(db *gorm.DB) *FindVoterQuery {
	return &FindVoterQuery{
		db: db,
	}
}

func (q *FindVoterQuery) Exec(where ...interface{}) (*entity.Voter, error) {
	var voter entity.Voter

	if err := q.db.First(&voter, where...).Error; err != nil {
		return nil, err
	}

	return &voter, nil
}
