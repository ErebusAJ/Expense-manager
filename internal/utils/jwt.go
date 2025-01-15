package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// GENERATE JWT
// Function to generate JWT auth token
func GenerateJWT(userID uuid.UUID, userRole string) (string, error){
	//Token claims
	claims := jwt.MapClaims{
		"user_id"	: userID,
		"user_role"	: userRole,
		"exp"		: time.Now().Add(24 * time.Hour).Unix(), //24hrs expiration time
		"iat"		: time.Now().Unix(), // Issued at time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Sign token with secret key
	godotenv.Load()

	key := os.Getenv("SECRET_KEY")
	if key == ""{
		return "", errors.New("error getting signed key")
	}

	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("error signing key")
	}

	return signedToken, nil
}