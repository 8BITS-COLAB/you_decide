package command

import (
	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CreatePartyCommand struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewCreatePartyCommand(db *gorm.DB, validator *validator.Validate) *CreatePartyCommand {
	return &CreatePartyCommand{
		db:        db,
		validator: validator,
	}
}

func (c *CreatePartyCommand) Exec(createPartyDTO dto.CreatePartyDTO) (*entity.Party, error) {
	party := entity.Party{
		Initials: createPartyDTO.Initials,
		Name:     createPartyDTO.Name,
	}

	if err := c.validator.Struct(party); err != nil {
		return nil, err
	}

	if err := c.db.Create(&party).Error; err != nil {
		return nil, err
	}

	return &party, nil

}
