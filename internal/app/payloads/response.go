package payloads

type LoginResponse struct {
	Token string `json:"token"`
}

type ViewCartResponse struct {
	ProductID       string `json:"product_id"`
	ProductName     string `json:"product_name"`
	ProductPrice    uint64 `json:"product_price"`
	ProductQuantity uint   `json:"product_quantity"`
	UpdatedAt       int64  `json:"-"`
}
