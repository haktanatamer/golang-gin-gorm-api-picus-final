package shopping_api

import (
	"api-gin/package/internal/api/product_api"
	"api-gin/package/internal/domain/shopping"
	"math"
	"time"

	"fmt"

	"api-gin/package/pkg/helper"
	"api-gin/package/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ShoppingController struct {
	shoppingService *shopping.ShoppingService
}

func NewShoppingController(service *shopping.ShoppingService) *ShoppingController {
	return &ShoppingController{
		shoppingService: service,
	}
}

// @Summary Get cart
// @Tags Cart
// @Accept  json
// @Produce  json
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Success 200 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /shopping/list [get]
func (sc *ShoppingController) GetCart(g *gin.Context) {

	userId := middleware.GetUserId(viper.GetString("server.secret"), g)

	cartList := sc.shoppingService.GetCartList(userId)

	cartPrice := sc.shoppingService.GetCartPrice(userId)

	var carts CartResponse
	carts.Total = cartPrice

	for _, p := range cartList {
		carts.Items = append(carts.Items, CartItemResponse{Name: p.Brand, Brand: p.Brand, Quantity: p.Quantity, Price: p.Price, ProductId: p.ProductId})
	}

	if carts.Total > 0 {
		g.JSON(http.StatusCreated, gin.H{
			"carts": carts,
		})
	} else {
		g.JSON(http.StatusCreated, gin.H{
			"status": "Basket is empty.",
		})
	}

}

// @Summary Add or update cart
// @Tags Cart
// @Accept  json
// @Produce  json
// @Param cartRequest body []CartRequest false "Cart Items"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Success 200 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /shopping/add [post]
func (sc *ShoppingController) AddToCart(g *gin.Context) {

	var reqs []CartRequest

	if err := g.ShouldBindJSON(&reqs); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error_message": helper.ErrRequestBody})
		g.Abort()
		return
	}

	userId := middleware.GetUserId(viper.GetString("server.secret"), g)
	var totalPrice float64

	for _, req := range reqs {

		httpCode, errDetail := helper.Valid(g, &req)

		if httpCode != http.StatusOK {
			g.JSON(httpCode, gin.H{"error_message": errDetail})
			g.Abort()
			return
		}

		product := sc.shoppingService.ExistProductById(req.ProductId)

		if product.Id == 0 {
			allProducts := sc.shoppingService.GetAllProducts()
			var resProducts []product_api.ResponseRequest
			for _, p := range allProducts {
				resProducts = append(resProducts, product_api.ResponseRequest{Name: p.Name, Brand: p.Brand, Id: int(p.Id), Sku: p.Sku, Category: p.Category.Name})
			}

			g.JSON(httpCode, gin.H{"error_message": fmt.Sprintf("Product did not found productId : %d, select product from the list.", req.ProductId), "Products": resProducts})
			g.Abort()
			return
		}

		if req.Quantity > product.Stock {
			g.JSON(httpCode, gin.H{"error_message": fmt.Sprintf("There is not enough stock for productId : %d, avaible stock : %d", req.ProductId, product.Stock)})
			g.Abort()
			return
		}

		err := sc.shoppingService.AddOrUpdateToCartItem(req.ProductId, req.Quantity, userId, product.Price)

		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error_message": "An error while adding cart."})
			g.Abort()
			return
		}
		totalPrice += (product.Price * float64(req.Quantity))
	}

	err := sc.shoppingService.AddOrUpdateToCart(userId, totalPrice)

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error_message": "An error while adding cart."})
		g.Abort()
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"status": "The products has been added/updated cart.",
	})

}

// @Summary Delete or update item to cart
// @Tags Cart
// @Accept  json
// @Produce  json
// @Param cartRequest body []CartRequest false "Delete Cart Items"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Success 200 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /shopping/delete [post]
func (sc *ShoppingController) DeleteToCart(g *gin.Context) {

	var reqs []CartRequest

	if err := g.ShouldBindJSON(&reqs); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error_message": helper.ErrRequestBody})
		g.Abort()
		return
	}

	userId := middleware.GetUserId(viper.GetString("server.secret"), g)
	totalPrice := sc.shoppingService.GetCartPrice(userId)

	for _, req := range reqs {

		httpCode, errDetail := helper.Valid(g, &req)

		if httpCode != http.StatusOK {
			g.JSON(httpCode, gin.H{"error_message": errDetail})
			g.Abort()
			return
		}

		cartItem := sc.shoppingService.ProductExistInCartItem(req.ProductId, userId)

		if cartItem.ProductId == 0 {
			cartItems := sc.shoppingService.GetCartList(userId)

			g.JSON(httpCode, gin.H{"error_message": fmt.Sprintf("Product did not found productId : %d, select product from the list.", req.ProductId), "Products": cartItems})
			g.Abort()
			return
		}
		deleteRow := false
		if req.Quantity > cartItem.Quantity {
			req.Quantity = cartItem.Quantity
			deleteRow = true
		}

		err := sc.shoppingService.DeleteToCartItem(req.ProductId, req.Quantity, userId, cartItem.Price, deleteRow)

		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error_message": "An error while adding cart."})
			g.Abort()
			return
		}
		totalPrice -= (cartItem.Price * float64(req.Quantity))

	}

	var err error
	if math.Round(totalPrice*100)/100 == 0 {

		err = sc.shoppingService.DeleteCartByUserId(userId)
	} else {

		err = sc.shoppingService.AddOrUpdateToCart(userId, totalPrice)
	}

	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error_message": "An error while adding cart."})
		g.Abort()
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"status": "The products has been deleted/updated cart.",
	})

}

// @Summary Create Order
// @Tags Order
// @Accept  json
// @Produce  json
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Success 200 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /order/add [post]
func (sc *ShoppingController) AddOrder(g *gin.Context) {
	userId := middleware.GetUserId(viper.GetString("server.secret"), g)

	cart := sc.shoppingService.GetCartByUserId(userId)

	if cart.Id == 0 {
		g.JSON(http.StatusInternalServerError, gin.H{"status": "Basket is empty."})
		g.Abort()
		return
	}

	cartItems := sc.shoppingService.GetCartItemByUserId(userId)

	order := sc.shoppingService.CreateOrder(cart, cartItems)

	g.JSON(http.StatusOK, gin.H{
		"status": fmt.Sprintf("The order was created. Order id : %d", order.Id),
	})
}

// @Summary Get All Order
// @Tags Order
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /order/list [get]
func (sc *ShoppingController) GetOrders(g *gin.Context) {

	userId := middleware.GetUserId(viper.GetString("server.secret"), g)

	order := sc.shoppingService.GetOrder(userId)

	g.JSON(http.StatusOK, gin.H{
		"orders": order,
	})
}

// @Summary Cancel Order
// @Tags Order
// @Accept  json
// @Produce  json
// @Param cancelOrderRequest body CancelOrderRequest false "Order Id"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Success 200 {object} map[string]string
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /order/cancel [post]
func (sc *ShoppingController) CancelOrders(g *gin.Context) {

	userId := middleware.GetUserId(viper.GetString("server.secret"), g)

	var req CancelOrderRequest

	httpCode, errDetail := helper.BindAndValid(g, &req)

	if httpCode != http.StatusOK {
		g.JSON(httpCode, gin.H{"error_message": errDetail})
		g.Abort()
		return
	}

	order := sc.shoppingService.GetOrderById(req.OrderId, userId)

	if order.Id == 0 {
		g.JSON(httpCode, gin.H{"error_message": "Order did not found."})
		g.Abort()
		return
	}

	a := helper.DateDiffDay(time.Now(), order.CreatedAt)

	if a > 14 {
		g.JSON(httpCode, gin.H{"error_message": "Order cannot be canceled."})
		g.Abort()
		return
	}

	err := sc.shoppingService.CancelOrderById(req.OrderId, userId)

	if err != nil {
		g.JSON(httpCode, gin.H{"error_message": "An error while canceled order."})
		g.Abort()
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"status": "Order has been canceled",
	})
}
