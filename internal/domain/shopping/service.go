package shopping

import (
	"api-gin/package/internal/domain/product"
)

type ShoppingService struct {
	r ShoppingRepository
}

func NewShoppingService(r ShoppingRepository) *ShoppingService {
	return &ShoppingService{
		r: r,
	}
}

func (ss *ShoppingService) ExistProductById(id int) product.Product {
	product := ss.r.ExistProductById(id)

	return product
}

func (ss *ShoppingService) GetAllProducts() []product.Product {
	products := ss.r.GetAllProducts()

	return products
}

func (ss *ShoppingService) AddOrUpdateToCartItem(productId, quantity, userId int, price float64) error {
	err := ss.r.AddOrUpdateToCartItem(productId, quantity, userId, price)

	return err
}

func (ss *ShoppingService) AddOrUpdateToCart(userId int, price float64) error {
	err := ss.r.AddOrUpdateToCart(userId, price)

	return err
}

func (ss *ShoppingService) GetCartList(userId int) []Cart_List {
	cartList := ss.r.GetCartList(userId)

	return cartList
}

func (ss *ShoppingService) GetCartPrice(userId int) float64 {
	price := ss.r.GetCartPrice(userId)

	return price
}

func (ss *ShoppingService) ProductExistInCartItem(id, userId int) Cart_Item {
	item := ss.r.ProductExistInCartItem(id, userId)

	return item
}

func (ss *ShoppingService) DeleteToCartItem(productId, quantity, userId int, price float64, deleteRow bool) error {
	err := ss.r.DeleteToCartItem(productId, quantity, userId, price, deleteRow)

	return err
}

func (ss *ShoppingService) DeleteCartByUserId(userId int) error {
	err := ss.r.DeleteCartByUserId(userId)

	return err
}

func (ss *ShoppingService) GetCartByUserId(userId int) Cart {
	cart := ss.r.GetCartByUserId(userId)

	return cart
}

func (ss *ShoppingService) CreateOrder(cart Cart, items []Cart_Item) Order {
	order := ss.r.CreateOrder(cart)
	if order.Id > 0 {
		for _, i := range items {
			_ = ss.r.CreateOrderItem(i, int(order.Id))
		}
	}

	return order
}

func (ss *ShoppingService) GetCartItemByUserId(userId int) []Cart_Item {
	cartItems := ss.r.GetCartItemByUserId(userId)

	return cartItems
}

func (ss *ShoppingService) GetOrder(userId int) []Order {
	order := ss.r.GetOrder(userId)

	return order
}

func (ss *ShoppingService) GetOrderById(id, userId int) Order {
	order := ss.r.GetOrderById(id, userId)

	return order
}

func (ss *ShoppingService) CancelOrderById(orderId, userId int) error {
	order := ss.GetOrderById(orderId, userId)

	err := ss.r.CancelOrderById(orderId)

	if err != nil {
		return err
	}

	for _, i := range order.OrderItems {
		err = ss.r.AddStock(i.ProductId, i.Quantity)
		if err != nil {
			return err
		}
	}

	return err
}
