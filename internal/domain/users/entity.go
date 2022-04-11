package users

import "time"

const (
	ADMIN_ROLE    int = 1
	CUSTOMER_ROLE int = 2
)

type User struct {
	Id        uint `gorm:"primaryKey"`
	Username  string
	Password  string
	CreatedAt time.Time `gorm:"<-:create"`
}

type User_Role struct {
	Id     uint `gorm:"primaryKey"`
	UserId int
	RoleId int
}

type Role struct {
	Id        uint `gorm:"primaryKey"`
	Role      string
	CreatedAt time.Time `gorm:"<-:create"`
}

type UserAndRoles struct {
	user  *User
	roles []string
}

type User_Token struct {
	UserId    int `gorm:"primaryKey"`
	Token     string
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func NewUser(username, password string) *User {
	return &User{
		Username: username,
		Password: password,
	}
}

func NewUserRole(userId int) *User_Role {
	return &User_Role{
		UserId: userId,
		RoleId: CUSTOMER_ROLE,
	}
}

func NewUserAdminRole(userId int) *User_Role {
	return &User_Role{
		UserId: userId,
		RoleId: ADMIN_ROLE,
	}
}
