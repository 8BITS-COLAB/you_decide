package entity

type Block struct {
	Index         int    `json:"index" gorm:"column:index;primaryKey"`
	Hash          string `json:"hash" gorm:"column:hash"`
	Nonce         string `json:"nonce" gorm:"column:nonce"`
	Party         *Party `json:"party" gorm:"foreignKey:PartyInitials"`
	PartyInitials string `json:"party_initials" gorm:"column:party_initials"`
	PrevHash      string `json:"prev_hash" gorm:"column:prev_hash"`
	Votes         []Vote `json:"votes" gorm:"foreignKey:BlockIndex"`
}
