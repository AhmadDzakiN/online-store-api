package payloads

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AddProductRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid4"`
	Quantity  uint   `json:"quantity" validate:"required,gt=0"`
}

type CheckoutRequest struct {
	CartID uint64                `json:"cart_id" validate:"required,numeric"`
	Items  []CheckoutItemRequest `json:"items" validate:"required"`
}

type CheckoutItemRequest struct {
	CartItemID uint64 `json:"cart_item_id" validate:"required,numeric"`
	ProductID  string `json:"product_id" validate:"required,uuid4"`
	Price      uint64 `json:"price" validate:"required,numeric,gt=0"`
	Quantity   uint   `json:"quantity" validate:"required,numeric,gt=0"`
}
