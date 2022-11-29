package command

import (
	"fmt"

	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/ElioenaiFerrari/youdecide/app/protocol"
	"github.com/ElioenaiFerrari/youdecide/app/query"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	valueobject "github.com/ElioenaiFerrari/youdecide/domain/value-object"
	"gorm.io/gorm"
)

type CreateCandidatureCommand struct {
	db                 *gorm.DB
	validator          protocol.ValidateStructProtocol
	findCandidateQuery *query.FindCandidateQuery
}

func NewCreateCandidatureCommand(db *gorm.DB, validator protocol.ValidateStructProtocol, findCandidateQuery *query.FindCandidateQuery) *CreateCandidatureCommand {
	return &CreateCandidatureCommand{
		db:                 db,
		validator:          validator,
		findCandidateQuery: findCandidateQuery,
	}
}

func (c *CreateCandidatureCommand) Exec(createCandidatureDTO dto.CreateCandidatureDTO) (*entity.Candidature, error) {
	candidate, err := c.findCandidateQuery.Exec("code = ?", createCandidatureDTO.CandidateCode)

	if err != nil {
		return nil, err
	}

	signature, _ := valueobject.NewSignature(fmt.Sprintf("%s:%s:%d", candidate.Code, createCandidatureDTO.Position, createCandidatureDTO.Year), 2)

	candidature := entity.Candidature{
		CandidateCode: candidate.Code,
		Position:      createCandidatureDTO.Position,
		Year:          createCandidatureDTO.Year,
		Signature:     signature,
	}

	if err := c.validator.ValidateStruct(candidature); err != nil {
		return nil, err
	}

	if err := c.db.Create(&candidature).Error; err != nil {
		return nil, err
	}

	return &candidature, nil

}
