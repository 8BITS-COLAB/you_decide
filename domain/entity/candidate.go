package entity

type Candidate struct {
	Candidatures  []Candidature `json:"candidatures" gorm:"foreignKey:CandidateCode"`
	Code          string        `json:"code" gorm:"column:code;primaryKey"`
	Name          string        `json:"name" gorm:"column:name"`
	Party         *Party        `json:"party" gorm:"foreignKey:PartyInitials"`
	PartyInitials string        `json:"party_initials" gorm:"column:party_initials"`
}
