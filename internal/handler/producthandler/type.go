package producthandler

type ProductResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	SalePrice   *float64 `json:"sale_price,omitempty"`
	Price       float64  `json:"price"`
}