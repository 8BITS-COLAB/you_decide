package dto

type CreateCandidatureDTO struct {
	CandidateName string `json:"candidate_name"`
	Code          string `json:"code"`
	Position      string `json:"position"`
	Year          int    `json:"year"`
}
