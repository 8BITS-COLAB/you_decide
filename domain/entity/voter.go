package entity

type Voter struct {
	Address string `json:"address" gorm:"column:address;primaryKey"`
	Email   string `json:"email" gorm:"column:email;uniqueIndex"`
	Name    string `json:"name" gorm:"column:name"`
	Phone   string `json:"phone" gorm:"column:phone;uniqueIndex"`
	Votes   []Vote `json:"votes" gorm:"foreignKey:VoterAddress"`
}
