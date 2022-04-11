package users_test

import (
	"api-gin/package/internal/domain/users"
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
	repository *users.UserRepository
	user       *users.User
}

func (s *Suite) SetupSuite() {

	db, err := database.TestSetup()

	if err != nil {
		return
	}
	s.db = db
	s.repository = users.NewUserRepository(s.db)

	for _, val := range getModels() {
		s.db.AutoMigrate(val)
	}
}

func getModels() []interface{} {
	return []interface{}{&users.User{}, &users.User_Role{}, &users.User_Token{}}
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
		tag  string
		user *users.User
	}{
		{"user1", users.NewUser("user1", "psw1")},
		{"user2", users.NewUser("user2", "psw2")},
		{"user3", users.NewUser("user3", "psw3")},
	}
	for _, test := range tests {
		err := s.repository.Create(test.user)
		assert.Equal(s.T(), nil, err, "Error should be nil")
	}

}

func (s *Suite) TestExistByUsername() {
	tests := []struct {
		tag  string
		user *users.User
		res  bool
	}{
		{"user1", users.NewUser("user1", "psw1"), true},
		{"user4", users.NewUser("user4", "psw4"), false},
	}
	for _, test := range tests {
		isUser := s.repository.ExistByUsername(test.user.Username)
		assert.Equal(s.T(), test.res, isUser, "")
	}

}

func (s *Suite) TestGetLoginUser() {
	tests := []struct {
		tag   string
		user  *users.User
		res   bool
		roles []string
	}{
		{"user2", users.NewUser("user2", "psw2"), false, []string{"customer"}},
	}
	for _, test := range tests {
		isUser, user, roles := s.repository.GetLoginUser(test.user.Username, test.user.Password)
		assert.Equal(s.T(), test.res, isUser, "")
		assert.Equal(s.T(), test.user, user, "")
		assert.Equal(s.T(), test.roles, roles, "")
	}

}

func (s *Suite) TestGetUserRoles() {
	tests := []struct {
		tag   string
		user  *users.User
		roles []string
	}{
		{"user1", &users.User{Id: 2}, []string{"customer"}},
	}
	for _, test := range tests {
		roles := s.repository.GetUserRoles(*test.user)
		assert.Equal(s.T(), test.roles, roles, "")
	}

}

func (s *Suite) TestAddTokenToUser() {
	tests := []struct {
		tag    string
		userId int
		token  string
		res    bool
	}{
		{"ex1", 2, "token", false},
		{"ex2", 10, "token2", false},
	}
	for _, test := range tests {
		res := s.repository.AddTokenToUser(test.userId, test.token)
		assert.Equal(s.T(), test.res, res, test.tag+" fail")
	}

}

func (s *Suite) TestGetUserToken() {
	tests := []struct {
		tag    string
		userId int
		token  string
	}{
		{"ex1", 2, "token"},
		{"ex2", 10, "token2"},
	}
	for _, test := range tests {
		res := s.repository.GetUserToken(test.userId)
		assert.Equal(s.T(), test.token, res, test.tag+" fail")
	}

}
