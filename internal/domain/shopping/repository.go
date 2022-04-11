package shopping

import (
	"api-gin/package/internal/domain/product"
	"math"

	"gorm.io/gorm"
)

type ShoppingRepository struct {
	db *gorm.DB
}

func NewShoppingRepository(db *gorm.DB) *ShoppingRepository {
	return &ShoppingRepository{
		db: db,
	}
}

func (r *ShoppingRepository) ExistProductById(id int) product.Product {
	var product product.Product
	r.db.Raw("SELECT * FROM product p where p.Id= ? and p.IsDeleted = ?", id, false).Scan(&product)

	return product
}

func (r *ShoppingRepository) GetAllProducts() []product.Product {
	var products []product.Product
	r.db.Preload("Category").Where("IsDeleted = ?", false).Find(&products)

	return products
}

func (r *ShoppingRepository) AddOrUpdateToCartItem(id, quantity, userId int, price float64) error {
	var cartItem Cart_Item
	r.db.Preload("Cart_Item").Where("CartId = ? and ProductId = ?", userId, id).Find(&cartItem)

	if cartItem.Id == 0 {
		cartItem.CartId = userId
		cartItem.ProductId = id
		cartItem.Quantity = quantity
		cartItem.Price = price
		if err := r.db.Create(&cartItem).Error; err != nil {
			return err
		}
	} else {
		cartItem.Quantity += quantity
		if err := r.db.Save(&cartItem).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *ShoppingRepository) AddOrUpdateToCart(userId int, price float64) error {

	var cart Cart
	r.db.Preload("Cart").Where("UserId = ? ", userId).Find(&cart)

	if cart.Id == 0 {
		cart.UserId = userId
		cart.Total = math.Round(price*100) / 100
		if err := r.db.Create(&cart).Error; err != nil {
			return err
		}
	} else {
		var cartTotal Cart_Total
		r.db.Raw("SELECT SUM(ROUND(i.Quantity*i.Price, 2)) as Price FROM basket_api.cart_item i where i.CartId = ?", userId).Scan(&cartTotal)

		cart.UserId = userId
		cart.Total = cartTotal.Price

		if err := r.db.Save(&cart).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *ShoppingRepository) GetCartList(userId int) []Cart_List {
	var cartList []Cart_List
	r.db.Raw("SELECT * FROM basket_api.v_cart_list v where v.CartId = ?", userId).Scan(&cartList)

	return cartList
}

func (r *ShoppingRepository) GetCartPrice(userId int) float64 {
	var cart Cart
	r.db.Preload("Cart").Where("UserId = ?", userId).Find(&cart)

	return cart.Total
}

func (r *ShoppingRepository) ProductExistInCartItem(id, userId int) Cart_Item {
	var cartItem Cart_Item
	r.db.Raw("SELECT * FROM v_cart_list v where v.ProductId = ? and v.CartId = ?", id, userId).Scan(&cartItem)

	return cartItem
}

func (r *ShoppingRepository) DeleteToCartItem(id, quantity, userId int, price float64, deleteRow bool) error {
	var cartItem Cart_Item
	r.db.Preload("Cart_Item").Where("CartId = ? and ProductId = ?", userId, id).Find(&cartItem)

	if deleteRow {
		r.db.Delete(cartItem)
	} else {
		cartItem.Quantity -= quantity
		if err := r.db.Save(&cartItem).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *ShoppingRepository) DeleteCartByUserId(userId int) error {

	var cart Cart
	r.db.Preload("Cart").Where("UserId = ? ", userId).Find(&cart)

	if cart.Id > 0 {
		r.db.Delete(cart)
	}

	return nil
}

func (r *ShoppingRepository) GetCartByUserId(userId int) Cart {
	var cart Cart
	r.db.Preload("Cart").Where("UserId = ?", userId).Find(&cart)

	return cart
}

func (r *ShoppingRepository) CreateOrder(cart Cart) Order {
	var order Order
	order.UserId = cart.UserId
	order.IsCanceled = false
	order.Total = cart.Total
	if err := r.db.Save(&order).Error; err != nil {
		return Order{}
	}
	r.db.Delete(cart)
	return order
}

func (r *ShoppingRepository) GetCartItemByUserId(userId int) []Cart_Item {
	var cartItems []Cart_Item
	r.db.Preload("Cart_Item").Where("CartId = ?", userId).Find(&cartItems)

	return cartItems
}

func (r *ShoppingRepository) CreateOrderItem(item Cart_Item, orderId int) error {
	var orderItem Order_Item
	orderItem.OrderId = orderId
	orderItem.Quantity = item.Quantity
	orderItem.ProductId = item.ProductId
	orderItem.Price = item.Price
	if err := r.db.Save(&orderItem).Error; err != nil {
		return err
	}
	r.db.Delete(item)
	return nil
}

func (r *ShoppingRepository) GetOrder(userId int) []Order {
	var order []Order
	r.db.Where("UserId = ?", userId).Preload("OrderItems").Preload("OrderItems.Products").Preload("OrderItems.Products.Category").Find(&order)

	return order
}

func (r *ShoppingRepository) GetOrderById(id, userId int) Order {
	var order Order
	r.db.Where("Id = ? and UserId = ?", id, userId).Preload("OrderItems").Find(&order)

	return order
}

func (r *ShoppingRepository) CancelOrderById(orderId int) error {
	var order Order
	r.db.Where("Id = ?", orderId).Find(&order)
	order.IsCanceled = true
	if err := r.db.Save(&order).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShoppingRepository) AddStock(productId, quantity int) error {
	var product product.Product
	r.db.Raw("SELECT * FROM product p where p.Id= ? and p.IsDeleted = ?", productId, false).Scan(&product)
	product.Stock += quantity
	if err := r.db.Save(&product).Error; err != nil {
		return err
	}

	return nil
}
