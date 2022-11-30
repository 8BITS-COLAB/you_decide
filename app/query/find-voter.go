package query

import (
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"gorm.io/gorm"
)

type FindCandidatureQuery struct {
	db *gorm.DB
}

func NewFindCandidatureQuery(db *gorm.DB) *FindCandidatureQuery {
	return &FindCandidatureQuery{
		db: db,
	}
}

func (q *FindCandidatureQuery) Exec(where ...interface{}) (*entity.Candidature, error) {
	var candidature entity.Candidature

	if err := q.db.First(&candidature, where...).Error; err != nil {
		return nil, err
	}

	return &candidature, nil
}
