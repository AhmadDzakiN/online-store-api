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
