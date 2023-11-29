package dto

type PagedDto struct {
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
	Value      any `json:"value"`
}
