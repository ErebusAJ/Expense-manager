package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// AUTHENTICATION FUNCTION
// Gets AuthHeader Bearer header for JWT Auth
// Sends userID as key value pair to the handlerfunc as context
func AuthMiddleware(tokenJWT string) gin.HandlerFunc{
	return func(c *gin.Context){
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer "){
			c.JSON(http.StatusUnauthorized, utils.MessageObj("authorization header missing"))
			log.Printf("unable to find auth header")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token * jwt.Token) (interface{}, error){
			return []byte(tokenJWT), nil
		} )
		if err != nil{
			c.JSON(http.StatusUnauthorized, utils.MessageObj("invalid token"))
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["user_id"].(string)

		userID, err := uuid.Parse(userId)
		if err != nil{
			c.IndentedJSON(http.StatusInternalServerError, utils.MessageObj("internal error"))
			log.Printf("error parsing user id to uuid %v \n", err)
		}
		c.Set("userID", userID)

		c.Next()
	}
}