package query

import (
	"errors"

	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"gorm.io/gorm"
)

type FindCandidateQuery struct {
	db *gorm.DB
}

func NewFindCandidateQuery(db *gorm.DB) *FindCandidateQuery {
	return &FindCandidateQuery{
		db: db,
	}
}

func (q *FindCandidateQuery) Exec(where ...interface{}) (*entity.Candidate, error) {
	var candidate entity.Candidate

	if err := q.db.First(&candidate, where...).Error; err != nil {
		return nil, errors.New("candidate not found")
	}

	return &candidate, nil
}
