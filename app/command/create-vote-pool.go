package command

import (
	"fmt"
	"time"

	"github.com/ElioenaiFerrari/youdecide/app/dto"
	valueobject "github.com/ElioenaiFerrari/youdecide/domain/value-object"
	"github.com/go-playground/validator/v10"
	"github.com/patrickmn/go-cache"
)

type CreateVotePoolCommand struct {
	validator *validator.Validate
	votesPool *cache.Cache
}

func NewCreateVotePoolCommand(
	validator *validator.Validate,
	votesPool *cache.Cache,
) *CreateVotePoolCommand {
	return &CreateVotePoolCommand{
		validator: validator,
		votesPool: votesPool,
	}
}

func (c *CreateVotePoolCommand) Exec(createVoteDTO dto.CreateVoteDTO) error {
	signature, _ := valueobject.NewSignature(createVoteDTO.Mnemonic, 3)
	key := fmt.Sprintf("%d:%s", time.Now().Unix(), signature)

	if err := c.votesPool.Add(key, createVoteDTO, cache.DefaultExpiration); err != nil {
		return err
	}

	return nil
}
