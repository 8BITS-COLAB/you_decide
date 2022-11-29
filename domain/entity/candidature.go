package entity

type Candidature struct {
	CandidateCode string     `json:"candidate_code" gorm:"column:candidate_code"`
	Candidate     *Candidate `json:"candidate" gorm:"foreignKey:CandidateCode"`
	Position      string     `json:"position" gorm:"column:position"`
	Signature     string     `json:"signature" gorm:"column:signature;primaryKey"`
	Year          int        `json:"year" gorm:"column:year"`
}

func NewCandidature(candidateCode, position string, year int) *Candidature {
	return &Candidature{}
}
