package command

import (
	"errors"
	"log"
	"time"

	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/ElioenaiFerrari/youdecide/app/query"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	valueobject "github.com/ElioenaiFerrari/youdecide/domain/value-object"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CreateVoteCommand struct {
	db                   *gorm.DB
	createBlockCommand   *CreateBlockCommand
	validator            *validator.Validate
	findCandidatureQuery *query.FindCandidatureQuery
	findVoterQuery       *query.FindVoterQuery
	findLastVoteQuery    *query.FindLastVoteQuery
}

func NewCreateVoteCommand(
	db *gorm.DB,
	createBlockCommand *CreateBlockCommand,
	validator *validator.Validate,
	findCandidatureQuery *query.FindCandidatureQuery,
	findVoterQuery *query.FindVoterQuery,
	findLastVoteQuery *query.FindLastVoteQuery,
) *CreateVoteCommand {
	return &CreateVoteCommand{
		db:                   db,
		createBlockCommand:   createBlockCommand,
		validator:            validator,
		findCandidatureQuery: findCandidatureQuery,
		findVoterQuery:       findVoterQuery,
		findLastVoteQuery:    findLastVoteQuery,
	}
}

func (c *CreateVoteCommand) Exec(createVoteDTO dto.CreateVoteDTO) error {
	candidature, err := c.findCandidatureQuery.Exec("code = ? AND year = ?", createVoteDTO.CandidatureCode, time.Now().Year())
	if err != nil {
		return err
	}

	voterAddress := valueobject.GetAddress(createVoteDTO.Mnemonic, createVoteDTO.Password)
	voter, err := c.findVoterQuery.Exec("address = ?", voterAddress)
	if err != nil {
		return err
	}

	vote := entity.Vote{
		CandidatureCode: candidature.Code,
		VoterAddress:    voter.Address,
	}

	if err := c.validator.Struct(vote); err != nil {
		return err
	}

	lastVote, err := c.findLastVoteQuery.Exec(voter.Address)

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		block, err := c.createBlockCommand.Exec([]entity.Vote{vote})
		if err != nil {
			return err
		}

		log.Printf("new block inserted with index: %d", block.Index)
		return nil
	}

	if lastVote.Candidature.Year == time.Now().Year() {
		return errors.New("voter already voted in this year")
	}

	block, err := c.createBlockCommand.Exec([]entity.Vote{vote})
	if err != nil {
		return err
	}

	log.Printf("new block inserted with index: %d", block.Index)

	return err
}
