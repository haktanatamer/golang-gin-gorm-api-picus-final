package shopping_api

type CartRequest struct {
	ProductId int `json:"productId" valid:"Match(^[0-9]*$); Required;Min(1)"`
	Quantity  int `json:"quantity" valid:"Match(^[0-9]*$); Required;Min(1);Max(7)"`
}

type CartResponse struct {
	Total float64            `json:"total"`
	Items []CartItemResponse `json:"items"`
}

type CartItemResponse struct {
	CartId    int     `json:"-"`
	ProductId int     `json:"id"`
	Name      string  `json:"name"`
	Brand     string  `json:"brand"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type CancelOrderRequest struct {
	OrderId int `json:"orderId"`
}
