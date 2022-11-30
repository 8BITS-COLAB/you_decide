package dto

type CreateVoteDTO struct {
	CandidatureCode string `json:"candidature_code"`
	Mnemonic        string `json:"mnemonic"`
	Password        string `json:"password"`
}
