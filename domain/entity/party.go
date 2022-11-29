package entity

type Party struct {
	Candidates []Candidate `json:"candidates" gorm:"foreignKey:PartyInitials"`
	Initials   string      `json:"initials" gorm:"column:initials;primaryKey"  validate:"required,min=2"`
	Name       string      `json:"name" gorm:"column:name"  validate:"required,min=3"`
}
