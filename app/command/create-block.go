package command

import (
	"time"

	"github.com/ElioenaiFerrari/youdecide/app/query"
	"github.com/ElioenaiFerrari/youdecide/domain/entity"
	valueobject "github.com/ElioenaiFerrari/youdecide/domain/value-object"
	"github.com/bytedance/sonic"
	"gorm.io/gorm"
)

type CreateBlockCommand struct {
	db                 *gorm.DB
	findLastBlockQuery *query.FindLastBlockQuery
}

func NewCreateBlockCommand(db *gorm.DB, findLastBlockQuery *query.FindLastBlockQuery) *CreateBlockCommand {
	return &CreateBlockCommand{
		db:                 db,
		findLastBlockQuery: findLastBlockQuery,
	}
}

func (c *CreateBlockCommand) Exec(votes []entity.Vote) (*entity.Block, error) {
	lastBlock, err := c.findLastBlockQuery.Exec()

	if err != nil {
		return nil, err
	}

	block := entity.Block{
		Votes:     votes,
		Index:     lastBlock.Index + 1,
		PrevHash:  lastBlock.Hash,
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
