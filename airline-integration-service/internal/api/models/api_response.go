package models

type ApiResponse struct {
	StatusCode    *int    `json:"status_code,omitempty"`
	StatusMessage string  `json:"status_message,omitempty"`
	Description   *string `json:"description,omitempty"`
	Data          *any    `json:"data,omitempty"`
	Took          *int64  `json:"took,omitempty"`
}

// func NewApiResponse(statusCode int, statusMessage string) *ApiResponse {
// 	return &ApiResponse{
// 		StatusCode: &statusCode,
// 		StatusMessage: statusMessage,
// 	}
// }

func NewErrorHandler(statusCode int, statusMessage string, description string) *ApiResponse {
	return &ApiResponse{
		StatusCode:    &statusCode,
		StatusMessage: statusMessage,
		Description:   &description,
	}
}

func NewApiResponse(statusCode int, statusMessage string, description string, data any, took int64) *ApiResponse {
	return &ApiResponse{
		StatusCode:    &statusCode,
		StatusMessage: statusMessage,
		Description:   &description,
		Data:          &data,
		Took:          &took,
	}
}
