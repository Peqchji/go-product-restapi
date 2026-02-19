package producthandler

type ProductResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	SalePrice   *float64 `json:"sale_price,omitempty"`
	Price       float64  `json:"price"`
}


// for swagger
type ProductSuccessResponse struct {
	Successful bool            `json:"successful"`
	ErrorCode  *string         `json:"error_code,omitempty"`
	Data       ProductResponse `json:"data"`
}
