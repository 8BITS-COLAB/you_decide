package command

import (
	"time"

	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	valueobject "github.com/ElioenaiFerrari/youdecide/domain/value-object"
	"github.com/bytedance/sonic"
	"gorm.io/gorm"
)

type CreateGenesisBlockCommand struct {
	db *gorm.DB
}

func NewCreateGenesisBlockCommand(db *gorm.DB) *CreateGenesisBlockCommand {
	return &CreateGenesisBlockCommand{
		db: db,
	}
}

func (c *CreateGenesisBlockCommand) Exec() (*entity.Block, error) {
	block := entity.Block{
		Votes:     []entity.Vote{},
		Index:     0,
		PrevHash:  "000",
		Timestamp: int(time.Now().UnixNano()),
	}

	blockBytes, _ := sonic.Marshal(block)
	signature, nonce := valueobject.NewSignature(string(blockBytes), 3)

	block.Hash = signature
	block.Nonce = nonce

	if err := c.db.Create(&block).Error; err != nil {
		return nil, err
	}

	return &block, nil
}
