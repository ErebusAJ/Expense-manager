package middleware

import (
	"net/http"
	"strings"

	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// ADMIN AUTH MIDDLEWARE
// Checks access level of incoming user provides access to only admin level access
func AdminMiddleware(signKey string) gin.HandlerFunc{
	return func (c *gin.Context){
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer "){
			utils.ErrorJSON(c, http.StatusUnauthorized, "invalid request", "error retrieving auth header from request", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
			return []byte(signKey), nil
		})
		if err != nil{
			utils.ErrorJSON(c, http.StatusUnauthorized, "invalid token", "invalid token", err)
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["user_id"]
		userRole := claims["user_role"]

		c.Set("userId", userId)
		c.Set("userRole", userRole)

		c.Next()
	}
}