package command

import (
	"github.com/ElioenaiFerrari/youdecide/app/dto"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	valueobject "github.com/ElioenaiFerrari/youdecide/domain/value-object"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CreateVoterCommand struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewCreateVoterCommand(db *gorm.DB, validator *validator.Validate) *CreateVoterCommand {
	return &CreateVoterCommand{
		db:        db,
		validator: validator,
	}
}

func (c *CreateVoterCommand) Exec(createVoterDTO dto.CreateVoterDTO) (*entity.Voter, *string, error) {
	address, mnemonic := valueobject.NewAddress(createVoterDTO.Password)

	voter := entity.Voter{
		Address: address,
		Email:   createVoterDTO.Email,
		Name:    createVoterDTO.Name,
		Phone:   createVoterDTO.Phone,
	}

	if err := c.validator.Struct(voter); err != nil {
		return nil, nil, err
	}

	if err := c.db.Create(&voter).Error; err != nil {
		return nil, nil, err
	}

	return &voter, &mnemonic, nil

}
