package entity

type Block struct {
	Index     int    `json:"index" gorm:"column:index;primaryKey" validate:"required,number"`
	Hash      string `json:"hash" gorm:"column:hash" validate:"required"`
	Nonce     int    `json:"nonce" gorm:"column:nonce" validate:"required,numner"`
	PrevHash  string `json:"prev_hash" gorm:"column:prev_hash" validate:"required"`
	Timestamp int    `json:"timestamp" gorm:"column:timestamp" validate:"required,number"`
	Votes     []Vote `json:"votes" gorm:"foreignKey:BlockIndex"`
}
