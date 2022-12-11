package entity

type Candidate struct {
	Candidatures  []Candidature `json:"candidatures" gorm:"foreignKey:CandidateName"`
	ImageURL      string        `json:"image_url" gorm:"column:image_url"`
	Name          string        `json:"name" gorm:"column:name;primaryKey" validate:"-"`
	Party         *Party        `json:"party" gorm:"foreignKey:PartyInitials"`
	PartyInitials string        `json:"party_initials" gorm:"column:party_initials" validate:"-"`
}
