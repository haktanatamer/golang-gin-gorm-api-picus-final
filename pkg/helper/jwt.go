package helper

import (
	"api-gin/package/internal/domain/users"
	jwtHelper "api-gin/package/pkg/jwt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

func CreateToken(user users.User, roles []string) string {

	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":   user.Id,
		"username": user.Username,
		"iat":      time.Now().Unix(),
		"iss":      os.Getenv("ENV"),
		"exp": time.Now().Add(100 *
			time.Minute).Unix(),
		"roles": roles,
	})

	jwtSecret := viper.GetString("server.secret")

	token := jwtHelper.GenerateToken(jwtClaims, jwtSecret)

	return token
}
