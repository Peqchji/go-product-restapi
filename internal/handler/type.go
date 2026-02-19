package handler


type BaseAPIResponse struct {
	Successful bool   `json:"successful"`
	ErrorCode  *string `json:"error_code,omitempty"`
	Data       any    `json:"data,omitempty"`
}