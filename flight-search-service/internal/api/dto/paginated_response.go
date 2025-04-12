package dto

type PaginatedResponse[T any] struct {
	Flights    []T `json:"flights"`
	Total      int `json:"total"`
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}
