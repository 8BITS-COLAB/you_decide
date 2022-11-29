package dto

type CreateCandidateDTO struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	PartyInitials string `json:"party_initials"`
}
