package model

type ErrorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
