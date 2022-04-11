package category

import (
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) ExistByName(c Category) bool {
	var category Category
	r.db.Where("Name = ?", c.Name).Find(&category)
	if category.Id > 0 {
		return true
	}

	return false
}

func (r *CategoryRepository) Create(newC *Category) error {
	if err := r.db.Create(&newC).Error; err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepository) GetAll() []Category {
	var categories []Category
	r.db.Where("Active = ?", true).Where("IsDeleted = ?", false).Find(&categories)

	return categories
}

func (r *CategoryRepository) GetAllByPagination(page, limit int) []Category {
	var categories []Category
	r.db.Limit(limit).Offset((page-1)*limit).Where("Active = ?", true).Where("IsDeleted = ?", false).Find(&categories)

	return categories
}
