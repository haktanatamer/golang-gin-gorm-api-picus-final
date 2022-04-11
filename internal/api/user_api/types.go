package user_api

type LoginRequest struct {
	Username string `json:"username" valid:"Required"`
	Password string `json:"password" valid:"Required"`
}
