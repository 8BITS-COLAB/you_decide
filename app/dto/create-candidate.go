package dto

type CreateCandidateDTO struct {
	ImageURL      string `json:"image_url"`
	Name          string `json:"name"`
	PartyInitials string `json:"party_initials"`
}
