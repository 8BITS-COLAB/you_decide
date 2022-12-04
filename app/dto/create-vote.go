package dto

type CreateVoteDTO struct {
	CandidatureCode string `json:"candidature_code"`
	Mnemonic        string `json:"mnemonic"`
	Email           string `json:"email"`
	Password        string `json:"password"`
}
