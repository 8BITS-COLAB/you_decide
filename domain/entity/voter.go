package entity

type Voter struct {
	Address string `json:"address" gorm:"column:address;primaryKey" validate:"-"`
	Email   string `json:"email" gorm:"column:email;uniqueIndex" validate:"required,email"`
	Name    string `json:"name" gorm:"column:name" validate:"required,min=3"`
	Phone   string `json:"phone" gorm:"column:phone;uniqueIndex" validate:"required,len=11"`
	Votes   []Vote `json:"votes" gorm:"foreignKey:VoterAddress"`
}
