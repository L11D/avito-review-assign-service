package dto

type ErrorDTO struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type FullErrorDTO struct {
	Error ErrorDTO `json:"error"`
}
