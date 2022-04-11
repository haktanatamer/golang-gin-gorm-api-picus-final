package product_test

import (
	"api-gin/package/internal/domain/product"
	"api-gin/package/pkg/database"
	"api-gin/package/pkg/helper"
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
	repository *product.ProductRepository
	product    *product.Product
}

func (s *Suite) SetupSuite() {

	db, err := database.TestSetup()

	if err != nil {
		return
	}
	s.db = db
	s.repository = product.NewProductRepository(s.db)

	for _, val := range getModels() {
		s.db.AutoMigrate(val)
	}
}

func getModels() []interface{} {
	return []interface{}{&product.Product{}}
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

func (s *Suite) TestCreate() {
	tests := []struct {
		tag string
		p   *product.Product
	}{
		{"ex1", product.NewProduct("name1", "brand1", 3)},
		{"ex2", product.NewProduct("name2", "brand2", 4)},
	}
	for _, test := range tests {
		err := s.repository.Create(test.p)
		assert.Equal(s.T(), nil, err, "Error should be nil")
	}

}

func (s *Suite) TestExistByName() {
	tests := []struct {
		tag string
		p   string
		res bool
	}{
		{"ex1", "name1", true},
		{"ex2", "name12", false},
	}
	for _, test := range tests {
		isUser := s.repository.ExistByName(test.p)
		assert.Equal(s.T(), test.res, isUser, "")
	}
}

func (s *Suite) TestGetAll() {
	res := s.repository.GetAll()
	assert.Equal(s.T(), 2, len(res), "Return values length should be 2")

}

func (s *Suite) TestGetAllByPagination() {
	res := s.repository.GetAllByPagination(1, 3)
	assert.Equal(s.T(), 2, len(res), "Return values length should be 3")

}

func (s *Suite) TestSearch() {
	res := s.repository.Search("name1")
	assert.Equal(s.T(), 1, len(res), "Return values length should be 1")

}

func (s *Suite) TestExistById() {

	tests := []struct {
		tag string
		id  int
		res bool
	}{
		{"ex1", 1, true},
		{"ex2", 1111, false},
	}
	for _, test := range tests {
		res := s.repository.ExistById(test.id)
		assert.Equal(s.T(), test.res, res, test.tag+"fail")
	}

}

func (s *Suite) TestDeleteById() {
	res := s.repository.DeleteById(2)
	assert.Equal(s.T(), nil, res, "Error should be nil")
}

func (s *Suite) TestUpdate() {
	tests := []struct {
		tag string
		p   *product.Product
		err error
	}{
		{"ex1", &product.Product{Id: 1, Name: "updatedName1", CategoryId: 3}, nil},
		{"ex2", &product.Product{Id: 111, CategoryId: 3222}, helper.ErrCategoryNotFound},
		{"ex3", &product.Product{Id: 1111, Name: "updatedName1", CategoryId: 3}, helper.ErrProductNotFound},
	}
	for _, test := range tests {
		res := s.repository.Update(test.p)
		assert.Equal(s.T(), test.err, res, test.tag+" fail")
	}
}
