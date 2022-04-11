package users

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) ExistByUsername(username string) bool {
	var user User
	r.db.Where("Username = ?", username).Find(&user)
	if user.Id > 0 {
		return true
	}

	return false
}

func (r *UserRepository) Create(newUser *User) error {
	if err := r.db.Create(&newUser).Error; err != nil {
		return err
	}

	err := r.AddRoleNewUser(int(newUser.Id))

	return err
}

func (r *UserRepository) AddRoleNewUser(userId int) error {
	newUserRole := NewUserRole(userId)
	if err := r.db.Create(&newUserRole).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetLoginUser(username, password string) (bool, User, []string) {
	var user User
	r.db.Where("Username = ?", username).Where("Password = ?", password).Find(&user)

	if user.Id == 0 {
		return true, user, nil
	}

	roles := r.GetUserRoles(user)

	return false, user, roles
}

func (r *UserRepository) GetUserRoles(u User) []string {
	var roles []string

	r.db.Raw("SELECT role FROM v_user_roles where id =?", u.Id).Scan(&roles)

	return roles
}

func (r *UserRepository) AddTokenToUser(userId int, token string) bool {
	var ut User_Token
	ut.Token = token
	ut.UserId = userId

	r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "UserId"}},
		DoUpdates: clause.AssignmentColumns([]string{"Token"}),
	}).Create(&ut)

	return false
}

func (r *UserRepository) GetUserToken(userId int) string {
	var ut User_Token
	r.db.Where("UserId = ?", userId).Find(&ut)
	if ut.UserId > 0 {
		return ut.Token
	}

	return ""
}
