package entity

type Candidate struct {
	Candidatures  []Candidature `json:"candidatures" gorm:"foreignKey:CandidateName"`
	Name          string        `json:"name" gorm:"column:name;primaryKey" validate:"required,min=3"`
	Party         *Party        `json:"party" gorm:"foreignKey:PartyInitials"`
	PartyInitials string        `json:"party_initials" gorm:"column:party_initials" validate:"required"`
}
