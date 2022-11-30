package command

import (
	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/ElioenaiFerrari/youdecide/app/query"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CreateCandidateCommand struct {
	db             *gorm.DB
	validator      *validator.Validate
	findPartyQuery *query.FindPartyQuery
}

func NewCreateCandidateCommand(db *gorm.DB, validator *validator.Validate, findPartyQuery *query.FindPartyQuery) *CreateCandidateCommand {
	return &CreateCandidateCommand{
		db:             db,
		validator:      validator,
		findPartyQuery: findPartyQuery,
	}
}

func (c *CreateCandidateCommand) Exec(createCandidateDTO dto.CreateCandidateDTO) (*entity.Candidate, error) {
	party, err := c.findPartyQuery.Exec("initials = ?", createCandidateDTO.PartyInitials)

	if err != nil {
		return nil, err
	}

	candidate := entity.Candidate{
		Name:          createCandidateDTO.Name,
		PartyInitials: party.Initials,
	}

	if err := c.validator.Struct(candidate); err != nil {
		return nil, err
	}

	if err := c.db.Create(&candidate).Error; err != nil {
		return nil, err
	}

	return &candidate, nil

}
