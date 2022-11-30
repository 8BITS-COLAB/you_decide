package query

import (
	"fmt"

	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"gorm.io/gorm"
)

type FindLastVoteQuery struct {
	db *gorm.DB
}

func NewFindLastVoteQuery(db *gorm.DB) *FindLastVoteQuery {
	return &FindLastVoteQuery{
		db: db,
	}
}

func (q *FindLastVoteQuery) Exec(voterAddress string) (*entity.Vote, error) {
	var vote entity.Vote

	if err := q.db.Preload("Candidature").Last(&vote, "voter_address = ?", voterAddress).Error; err != nil {
		return nil, err
	}

	fmt.Printf("%+v", vote)

	return &vote, nil
}
