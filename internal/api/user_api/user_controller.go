package user_api

import (
	"api-gin/package/internal/domain/users"
	"api-gin/package/pkg/helper"
	jwtHelper "api-gin/package/pkg/jwt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type UserController struct {
	userService *users.UserService
}

func NewUserController(service *users.UserService) *UserController {
	return &UserController{
		userService: service,
	}
}

// @Summary Add user
// @Tags User
// @Accept  json
// @Produce  json
// @Param loginRequest body LoginRequest false "User Info"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Success 200 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /user/add [post]
func (uc *UserController) Add(g *gin.Context) {

	var req LoginRequest

	httpCode, errDetail := helper.BindAndValid(g, &req)

	if httpCode != http.StatusOK {
		g.JSON(httpCode, gin.H{"error_message": errDetail})
		g.Abort()
		return
	}

	hasUser := uc.userService.ExistByUsername(req.Username)

	if hasUser {
		g.JSON(http.StatusConflict, gin.H{
			"error_message": "Username Already Exists.",
		})
		g.Abort()
		return
	}

	newUser := users.NewUser(req.Username, req.Password)

	err := uc.userService.Create(newUser)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "An error while adding user.",
		})
		g.Abort()
		return
	}

	g.JSON(http.StatusCreated, gin.H{
		"status": "The user has been created. Username : " + req.Username + " You can login to the system.",
	})
}

// @Summary Login user
// @Tags User
// @Accept  json
// @Produce  json
// @Param loginRequest body LoginRequest false "User Info"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/login [post]
func (uc *UserController) Login(g *gin.Context) {

	var req LoginRequest

	httpCode, errDetail := helper.BindAndValid(g, &req)

	if httpCode != http.StatusOK {
		g.JSON(httpCode, gin.H{"error_message": errDetail})
		g.Abort()
		return
	}

	hata, user, roles := uc.userService.GetLoginUser(req.Username, req.Password)

	if hata {
		g.JSON(http.StatusNotFound, gin.H{
			"error_message": "User not found!",
		})
		g.Abort()
		return
	}

	token := uc.userService.GetUserToken(int(user.Id))

	if token == "" {
		token = helper.CreateToken(user, roles)
	} else {
		_, err := jwtHelper.VerifyToken(token, viper.GetString("server.secret"), os.Getenv("ENV"))
		if err != nil {
			token = helper.CreateToken(user, roles)
		}
	}

	err := uc.userService.AddTokenToUser(int(user.Id), token)
	if err {
		g.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "An error while creating token.",
		})
		g.Abort()
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"status": "Login completed.", "token": token,
	})
}
