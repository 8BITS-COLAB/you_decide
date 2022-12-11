package query

import (
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"gorm.io/gorm"
)

type ListCandidaturesQuery struct {
	db *gorm.DB
}

func NewListCandidaturesQuery(db *gorm.DB) *ListCandidaturesQuery {
	return &ListCandidaturesQuery{
		db: db,
	}
}

func (q *ListCandidaturesQuery) Exec(where ...interface{}) ([]entity.Candidature, error) {
	var candidatures []entity.Candidature

	if err := q.db.Preload("Candidate").Find(&candidatures, where...).Error; err != nil {
		return nil, err
	}

	return candidatures, nil
}
