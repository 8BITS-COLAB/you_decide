package entity

type Vote struct {
	Block           *Block       `json:"block" gorm:"foreignKey:BlockIndex"`
	BlockIndex      int          `json:"block_index" gorm:"column:block_index" validate:"-"`
	Candidature     *Candidature `json:"candidature" gorm:"foreignKey:CandidatureCode"`
	CandidatureCode string       `json:"candidature_code" gorm:"column:candidature_code" validate:"required"`
	Voter           *Voter       `json:"voter" gorm:"foreignKey:VoterAddress"`
	VoterAddress    string       `json:"voter_address" gorm:"column:voter_address" validate:"required"`
}
