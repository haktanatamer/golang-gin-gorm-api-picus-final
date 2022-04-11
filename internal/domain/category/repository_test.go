package category_test

import (
	"api-gin/package/internal/domain/category"
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
	repository *category.CategoryRepository
	user       *category.Category
}

func (s *Suite) SetupSuite() {

	db, err := database.TestSetup()

	if err != nil {
		return
	}
	s.db = db
	s.repository = category.NewCategoryRepository(s.db)

	for _, val := range getModels() {
		s.db.AutoMigrate(val)
	}
}

func getModels() []interface{} {
	return []interface{}{&category.Category{}}
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
		c   *category.Category
	}{
		{"ex1", &category.Category{Name: "cate1", Active: true}},
		{"ex2", &category.Category{Name: "cate2", Active: true}},
	}
	for _, test := range tests {
		err := s.repository.Create(test.c)
		assert.Equal(s.T(), nil, err, "Error should be nil")
	}

}

func (s *Suite) TestGetAll() {
	res := s.repository.GetAll()
	assert.Equal(s.T(), 2, len(res), "Return values length should be 2")

}

func (s *Suite) TestGetAllByPagination() {
	res := s.repository.GetAllByPagination(1, 3)
	assert.Equal(s.T(), 3, len(res), "Return values length should be 3")

}
