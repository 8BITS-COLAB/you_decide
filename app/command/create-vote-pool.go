package command

import (
	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/go-playground/validator/v10"
)

type CreateVotePoolCommand struct {
	validator *validator.Validate
	votesPool *[]dto.CreateVoteDTO
}

func NewCreateVotePoolCommand(
	validator *validator.Validate,
	votesPool *[]dto.CreateVoteDTO,
) *CreateVotePoolCommand {
	return &CreateVotePoolCommand{
		validator: validator,
		votesPool: votesPool,
	}
}

func (c *CreateVotePoolCommand) Exec(createVoteDTO dto.CreateVoteDTO) error {
	*c.votesPool = append(*c.votesPool, createVoteDTO)

	return nil
}
