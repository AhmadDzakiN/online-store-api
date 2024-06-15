package payloads

type LoginResponse struct {
	Token string `json:"token"`
}

type ViewCartResponse struct {
	CartID uint64             `json:"cart_id"`
	Items  []CartItemResponse `json:"items"`
}

type CartItemResponse struct {
	CartItemID      uint64 `json:"cart_item_id"`
	ProductID       string `json:"product_id"`
	ProductName     string `json:"product_name"`
	ProductPrice    uint64 `json:"product_price"`
	ProductQuantity uint   `json:"product_quantity"`
	UpdatedAt       int64  `json:"-"` //For pagination
}

type ViewProductResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Price     uint64 `json:"price"`
	UpdatedAt int64  `json:"-"` //For pagination
}

type CheckoutResponse struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	TotalAmount uint64 `json:"total_amount"`
}
