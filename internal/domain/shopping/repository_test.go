package shopping_test

import (
	"api-gin/package/internal/domain/shopping"
	"api-gin/package/pkg/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
	db         *gorm.DB
	repository *shopping.ShoppingRepository
	cart       *shopping.Cart
}

func (s *Suite) SetupSuite() {

	db, err := database.TestSetup()

	if err != nil {
		return
	}
	s.db = db
	s.repository = shopping.NewShoppingRepository(s.db)

	for _, val := range getModels() {
		s.db.AutoMigrate(val)
	}
}

func getModels() []interface{} {
	return []interface{}{&shopping.Cart{}, &shopping.Cart_Item{}, &shopping.Order{}, &shopping.Order_Item{}}
}

/*
func (t *Suite) TearDownSuite() {
	sqlDB, _ := t.db.DB()
	defer sqlDB.Close()

	// Drop Table
	for _, val := range getModels() {
		t.db.Migrator().DropTable(val)
	}
}
*/

func (s *Suite) TestAddOrUpdateToCartItem() {
	tests := []struct {
		tag      string
		id       int
		quantity int
		userId   int
		price    float64
	}{
		{"ex1", 1, 2, 1, 15.25},
		{"ex2", 2, 1, 4, 58.62},
	}
	for _, test := range tests {
		res := s.repository.AddOrUpdateToCartItem(test.id, test.quantity, test.userId, test.price)
		assert.Equal(s.T(), nil, res, " fail")
	}

}

func (s *Suite) TestAddOrUpdateToCart() {
	tests := []struct {
		tag    string
		userId int
		price  float64
	}{
		{"ex1", 1, 42.45},
		{"ex2", 4, 77.64},
	}
	for _, test := range tests {
		res := s.repository.AddOrUpdateToCart(test.userId, test.price)
		assert.Equal(s.T(), nil, res, " fail")
	}

}

func (s *Suite) TestGetCartList() {
	tests := []struct {
		tag           string
		userId        int
		cartItemCount int
	}{
		{"ex1", 1, 3},
		{"ex2", 4, 2},
	}
	for _, test := range tests {
		res := s.repository.GetCartList(test.userId)
		assert.Equal(s.T(), test.cartItemCount, len(res), test.tag+" fail")
	}
}

func (s *Suite) TestGetCartPrice() {
	tests := []struct {
		tag       string
		userId    int
		cartPrice float64
	}{
		{"ex1", 1, 105.65},
		{"ex2", 4, 61.80},
	}
	for _, test := range tests {
		res := s.repository.GetCartPrice(test.userId)
		assert.Equal(s.T(), test.cartPrice, res, test.tag+" fail")
	}
}

func (s *Suite) TestProductExistInCartItem() {
	tests := []struct {
		tag       string
		productId int
		userId    int
	}{
		{"ex1", 1, 1}, {"ex2", 2, 4},
	}
	for _, test := range tests {
		res := s.repository.ProductExistInCartItem(test.productId, test.userId)
		assert.Equal(s.T(), test.userId, res.CartId, test.tag+" fail")
	}
}

func (s *Suite) TestDeleteToCartItem() {
	tests := []struct {
		tag       string
		id        int
		quantity  int
		userId    int
		price     float64
		deleteRow bool
	}{
		{"ex1", 1, 4, 1, 15.25, true}, {"ex2", 2, 1, 4, 58.62, false},
	}
	for _, test := range tests {
		res := s.repository.DeleteToCartItem(test.id, test.quantity, test.userId, test.price, test.deleteRow)
		assert.Equal(s.T(), nil, res, test.tag+" fail")
	}
}

func (s *Suite) TestDeleteCartByUserId() {
	tests := []struct {
		tag    string
		userId int
	}{
		{"ex1", 1}, {"ex2", 4},
	}
	for _, test := range tests {
		res := s.repository.DeleteCartByUserId(test.userId)
		assert.Equal(s.T(), nil, res, test.tag+" fail")
	}
}

func (s *Suite) TestGetCartByUserId() {
	tests := []struct {
		tag    string
		userId int
	}{
		{"ex1", 1}, {"ex2", 4},
	}
	for _, test := range tests {
		res := s.repository.GetCartByUserId(test.userId)
		assert.Equal(s.T(), test.userId, res.UserId, test.tag+" fail")
	}
}

func (s *Suite) TestCreateOrder() {
	tests := []struct {
		tag string
		c   *shopping.Cart
	}{
		{"ex1", &shopping.Cart{Id: 1, UserId: 1, Total: 42.45}},
		{"ex2", &shopping.Cart{Id: 2, UserId: 4, Total: 77.64}},
	}
	for _, test := range tests {
		res := s.repository.CreateOrder(*test.c)
		assert.Equal(s.T(), test.c.UserId, res.UserId, test.tag+" fail")
	}

}

func (s *Suite) TestGetCartItemByUserId() {
	tests := []struct {
		tag       string
		userId    int
		itemCount int
	}{
		{"ex1", 1, 3},
		{"ex2", 4, 2},
	}
	for _, test := range tests {
		res := s.repository.GetCartItemByUserId(test.userId)
		assert.Equal(s.T(), test.itemCount, len(res), test.tag+" fail")
	}

}

func (s *Suite) TestCreateOrderItem() {
	tests := []struct {
		tag     string
		ci      *shopping.Cart_Item
		orderId int
	}{
		{"ex1", &shopping.Cart_Item{CartId: 1, ProductId: 1, Quantity: 4, Price: 15.25}, 1},
		{"ex2", &shopping.Cart_Item{CartId: 4, ProductId: 2, Quantity: 1, Price: 58.62}, 2},
	}
	for _, test := range tests {
		res := s.repository.CreateOrderItem(*test.ci, test.orderId)
		assert.Equal(s.T(), nil, res, "Error should be nil")
	}
}

func (s *Suite) TestGetOrder() {
	tests := []struct {
		tag        string
		userId     int
		orderCount int
	}{
		{"ex1", 1, 2},
		{"ex2", 4, 1},
	}
	for _, test := range tests {
		res := s.repository.GetOrder(test.userId)
		assert.Equal(s.T(), test.orderCount, len(res), test.tag+" fail")
	}
}

func (s *Suite) TestGetOrderById() {
	tests := []struct {
		tag     string
		userId  int
		orderId int
	}{
		{"ex1", 1, 1},
		{"ex2", 4, 2},
	}
	for _, test := range tests {
		res := s.repository.GetOrderById(test.orderId, test.userId)
		assert.Equal(s.T(), test.orderId, res.Id, test.tag+" fail")
	}
}

func (s *Suite) TestCancelOrderById() {
	tests := []struct {
		tag     string
		orderId int
	}{
		{"ex1", 1},
		{"ex2", 2},
	}
	for _, test := range tests {
		res := s.repository.CancelOrderById(test.orderId)
		assert.Equal(s.T(), nil, res, "Error should be nil")
	}
}

func (s *Suite) TestAddStock() {
	tests := []struct {
		tag       string
		productId int
		quantity  int
	}{
		{"ex1", 1, 3},
		{"ex2", 2, 5},
	}
	for _, test := range tests {
		res := s.repository.AddStock(test.productId, test.quantity)
		assert.Equal(s.T(), nil, res, "Error should be nil")
	}
}
