package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)
func AuthMiddleware( jwtSecret string ) gin.HandlerFunc{
	return func( c *gin.Context){
		// 1. Authorization headerrrr
		authHeader := c.GetHeader("Authorization")

			if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})
			c.Abort()
			return
		}

		part:= strings.SplitN(authHeader , " ",2)

		if len(part) != 2 || part[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			c.Abort()
			return
		}

		tokenString := part[1]
		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}

return []byte(jwtSecret), nil
			},
		)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}
claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(string)

		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "user id missing from token",
			})
			c.Abort()
			return
		}

			c.Set("user_id", userID)

			c.Next()
	}
}