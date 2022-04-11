package users

type UserService struct {
	r UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{
		r: r,
	}
}

func (us *UserService) ExistByUsername(username string) bool {
	hasUser := us.r.ExistByUsername(username)

	return hasUser
}

func (us *UserService) Create(newUser *User) error {
	err := us.r.Create(newUser)

	return err
}

func (us *UserService) GetLoginUser(username, password string) (bool, User, []string) {
	hata, user, roles := us.r.GetLoginUser(username, password)

	return hata, user, roles
}

func (us *UserService) AddTokenToUser(userId int, token string) bool {
	err := us.r.AddTokenToUser(userId, token)

	return err
}

func (us *UserService) GetUserToken(userId int) string {
	token := us.r.GetUserToken(userId)

	return token
}
