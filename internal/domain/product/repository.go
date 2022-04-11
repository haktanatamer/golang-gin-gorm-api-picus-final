package product

import (
	"api-gin/package/internal/domain/category"
	"api-gin/package/pkg/helper"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) GetCategoryById(categoryId int) bool {
	var category category.Category
	r.db.Raw("SELECT * FROM category c where c.id= ?", categoryId).Scan(&category)
	if category.Id > 0 {
		return true
	}

	return false
}

func (r *ProductRepository) GetAllCategories() []category.Category {
	var categories []category.Category
	r.db.Raw("SELECT * FROM category ").Scan(&categories)

	return categories
}

func (r *ProductRepository) ExistByName(name string) bool {
	var product Product
	r.db.Where("Name = ?", name).Find(&product)
	if product.Id > 0 {
		return true
	}

	return false
}

func (r *ProductRepository) Create(newProduct *Product) error {
	if err := r.db.Create(&newProduct).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetAll() []Product {
	var products []Product
	r.db.Preload("Category").Where("IsDeleted = ?", false).Find(&products)

	return products
}

func (r *ProductRepository) Search(value string) []Product {
	var products []Product
	r.db.Joins("Product").Joins("Category").Where(
		r.db.Where("Product.IsDeleted = ?", false).Where(r.db.Where("Product.Name LIKE ?", "%"+value+"%").Or("Product.Sku LIKE ?", "%"+value+"%").Or("Product.Brand LIKE ?", "%"+value+"%").Or("Category.Name LIKE ?", "%"+value+"%")),
	).Find(&products)
	return products
}

func (r *ProductRepository) ExistById(id int) bool {
	var product Product
	r.db.Where("Id = ?", id).Find(&product)
	if product.Id > 0 {
		return true
	}

	return false
}

func (r *ProductRepository) DeleteById(id int) error {
	var product Product
	r.db.First(&product, id)
	product.IsDeleted = true
	if err := r.db.Save(&product).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) Update(p *Product) error {
	var product Product

	if p.CategoryId > 0 {
		if !r.GetCategoryById(p.CategoryId) {
			return helper.ErrCategoryNotFound
		}
		product.CategoryId = p.CategoryId
	}

	r.db.First(&product, p.Id)

	if product.Id == 0 {
		return helper.ErrProductNotFound
	}

	if len(p.Brand) > 0 {
		product.Brand = p.Brand
	}

	if len(p.Name) > 0 {
		product.Name = p.Name
	}

	if err := r.db.Save(&product).Error; err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) GetAllByPagination(page, limit int) []Product {
	var products []Product
	r.db.Limit(limit).Offset((page-1)*limit).Where("IsDeleted = ?", false).Find(&products)

	return products
}
