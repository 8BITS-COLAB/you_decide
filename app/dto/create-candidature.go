package dto

type CreateCandidatureDTO struct {
	CandidateCode string `json:"candidate_code"`
	Position      string `json:"position"`
	Year          int    `json:"year"`
}
