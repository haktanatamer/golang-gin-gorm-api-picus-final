package product

import (
	"api-gin/package/internal/domain/category"
	"api-gin/package/pkg/helper"

	"time"

	"bytes"
	"strconv"
)

type Product struct {
	Id                uint `gorm:"primaryKey"`
	Name, Brand, Sku  string
	CategoryId, Stock int
	IsDeleted         bool
	Price             float64
	Category          category.Category `gorm:"foreignKey:CategoryId;references:Id"`
	CreatedAt         time.Time         `gorm:"<-:create" json:"-"`
	UpdatedAt         time.Time         `json:"-"`
}

func NewProduct(name string, brand string, categoryId int) *Product {
	sku := createSku(name)
	return &Product{
		Name:       name,
		Brand:      brand,
		Sku:        sku,
		CategoryId: categoryId,
		Price:      helper.RandomFloatCreator(),
		IsDeleted:  false,
		Stock:      helper.RandomIntegerCreator(0, 100),
	}
}

// createSku generate random sku for product
func createSku(name string) string {
	first2 := name[0:2]
	var b bytes.Buffer
	b.WriteString(helper.StringUpper(first2))
	b.WriteString("-")
	b.WriteString(strconv.Itoa(helper.RandomIntegerCreator(100000, 999999)))

	return b.String()
}
