package middleware

import (
	"net/http"
	"os"

	jwtHelper "api-gin/package/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware authorization check
func AuthMiddleware(secretKey string, r string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.GetHeader("Authorization") != "" {
			decodedClaims, err := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey, os.Getenv("ENV"))
			if err != nil && err.Error() == "Token is expired" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Token is expired!"})
				c.Abort()
				return
			}

			if decodedClaims != nil {
				for _, role := range decodedClaims.Roles {
					if role == r || role == "admin" {
						c.Next()
						c.Abort()
						return
					}
				}
			}

			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to use this endpoint!"})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized!"})
		}
		c.Abort()
		return
	}
}

// GetUserId get userId from token
func GetUserId(secretKey string, c *gin.Context) int {
	decodedClaims, _ := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey, os.Getenv("ENV"))

	return decodedClaims.UserId

}
