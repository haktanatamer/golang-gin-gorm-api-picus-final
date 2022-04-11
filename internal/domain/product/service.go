package product

import (
	"api-gin/package/internal/domain/category"
)

type ProductService struct {
	r ProductRepository
}

func NewProductService(r ProductRepository) *ProductService {
	return &ProductService{
		r: r,
	}
}

func (ps *ProductService) GetCategoryById(categoryId int) bool {
	hasCategory := ps.r.GetCategoryById(categoryId)

	return hasCategory
}

func (ps *ProductService) GetAllCategories() []category.Category {
	allCategories := ps.r.GetAllCategories()

	return allCategories
}

func (ps *ProductService) ExistByName(name string) bool {
	hasProduct := ps.r.ExistByName(name)

	return hasProduct
}

func (ps *ProductService) Create(newProduct *Product) error {
	err := ps.r.Create(newProduct)

	return err
}

func (ps *ProductService) GetAll() []Product {
	result := ps.r.GetAll()

	return result
}

func (ps *ProductService) Search(value string) []Product {
	result := ps.r.Search(value)

	return result
}

func (ps *ProductService) ExistById(id int) bool {
	hasProduct := ps.r.ExistById(id)

	return hasProduct
}

func (ps *ProductService) DeleteById(id int) error {
	err := ps.r.DeleteById(id)

	return err
}

func (ps *ProductService) Update(p *Product) error {
	err := ps.r.Update(p)

	return err
}

func (ps *ProductService) GetAllByPagination(page, limit int) []Product {
	result := ps.r.GetAllByPagination(page, limit)

	return result
}
