package entity

type Vote struct {
	Block         *Block     `json:"block" gorm:"foreignKey:BlockIndex"`
	BlockIndex    int        `json:"block_index" gorm:"column:block_index"`
	Candidate     *Candidate `json:"candidate" gorm:"foreignKey:CandidateCode"`
	CandidateCode string     `json:"candidate_code" gorm:"column:candidate_code"`
	Email         string     `json:"email" gorm:"column:email;uniqueIndex"`
	Name          string     `json:"name" gorm:"column:name"`
	Phone         string     `json:"phone" gorm:"column:phone;uniqueIndex"`
	Signature     string     `json:"signature" gorm:"column:signature;primaryKey"`
	Voter         *Voter     `json:"voter" gorm:"foreignKey:VoterAddress"`
	VoterAddress  string     `json:"voter_address" gorm:"column:voter_address"`
}
