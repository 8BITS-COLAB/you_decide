package entity

type Candidature struct {
	Code          string     `json:"code" gorm:"column:code"`
	CandidateName string     `json:"candidate_name" gorm:"column:candidate_name"`
	Candidate     *Candidate `json:"candidate" gorm:"foreignKey:CandidateName"`
	Position      string     `json:"position" gorm:"column:position"`
	Signature     string     `json:"signature" gorm:"column:signature;primaryKey"`
	Year          int        `json:"year" gorm:"column:year"`
}

func NewCandidature(candidateCode, position string, year int) *Candidature {
	return &Candidature{}
}
