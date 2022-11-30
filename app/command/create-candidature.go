package command

import (
	"fmt"

	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/ElioenaiFerrari/youdecide/app/query"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	valueobject "github.com/ElioenaiFerrari/youdecide/domain/value-object"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CreateCandidatureCommand struct {
	db                 *gorm.DB
	validator          *validator.Validate
	findCandidateQuery *query.FindCandidateQuery
}

func NewCreateCandidatureCommand(db *gorm.DB, validator *validator.Validate, findCandidateQuery *query.FindCandidateQuery) *CreateCandidatureCommand {
	return &CreateCandidatureCommand{
		db:                 db,
		validator:          validator,
		findCandidateQuery: findCandidateQuery,
	}
}

func (c *CreateCandidatureCommand) Exec(createCandidatureDTO dto.CreateCandidatureDTO) (*entity.Candidature, error) {
	candidate, err := c.findCandidateQuery.Exec("name = ?", createCandidatureDTO.CandidateName)

	if err != nil {
		return nil, err
	}

	signature, _ := valueobject.NewSignature(fmt.Sprintf("%s:%d", createCandidatureDTO.Position, createCandidatureDTO.Year), 3)

	candidature := entity.Candidature{
		CandidateName: candidate.Name,
		Position:      createCandidatureDTO.Position,
		Year:          createCandidatureDTO.Year,
		Signature:     signature,
		Code:          createCandidatureDTO.Code,
	}

	if err := c.validator.Struct(candidature); err != nil {
		return nil, err
	}

	if err := c.db.Create(&candidature).Error; err != nil {
		return nil, err
	}

	return &candidature, nil

}
