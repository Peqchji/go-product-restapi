package handler

type BaseAPIResponse struct {
	Successful bool   `json:"successful"`
	ErrorCode  string `json:"error_code"`
	Data       *any   `json:"data,omitempty"`
}