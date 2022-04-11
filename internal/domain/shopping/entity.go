package shopping

import (
	"api-gin/package/internal/domain/product"
	"time"
)

type Cart struct {
	Id        uint `gorm:"primaryKey"`
	UserId    int
	Total     float64
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type Cart_Item struct {
	Id        uint `gorm:"primaryKey"`
	CartId    int
	ProductId int
	Quantity  int
	Price     float64
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type Order struct {
	Id         uint `gorm:"primaryKey" json:"orderId"`
	UserId     int  `json:"-"`
	Total      float64
	IsCanceled bool
	CreatedAt  time.Time    `json:"orderDate"`
	UpdatedAt  time.Time    `json:"-"`
	OrderItems []Order_Item `gorm:"foreignKey:OrderId;references:Id"`
}

type Order_Item struct {
	Id        uint `gorm:"primaryKey" json:"-"`
	OrderId   int  `json:"-"`
	ProductId int
	Quantity  int
	Price     float64
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	Products  product.Product `gorm:"foreignKey:ProductId;references:Id"`
}

type Cart_Total struct {
	Price float64
}

type Cart_List struct {
	ProductId int     `json:"id"`
	CartId    int     `json:"-"`
	Name      string  `json:"name"`
	Brand     string  `json:"brand"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"total_price"`
}
