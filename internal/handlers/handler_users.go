package handlers

import (
	"log"
	"net/http"

	"github.com/ErebusAJ/expense-manager/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CREATING NEW USER
// hashPassword function is used to hash the password recieved
// addUser creates the uses using sqlc DB queries

func hashPassword(password string)(string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}

	return string(hashedPassword), nil
}


func(cfg *apiConfig) addUser(c *gin.Context){
	type user struct{
		Id 			uuid.UUID	`json:"id"`
		Name		string		`json:"name"`
		Email		string		`json:"email"`
		Password	string		`json:"password"`
	}

	var param user
	c.BindJSON(&param)

	//Hashing Password
	hashPass, err := hashPassword(param.Password)
	if err != nil{
		log.Printf("error hashing password: %v", hashPass)
		return
	}

	err = cfg.DB.CreateUser(c, db.CreateUserParams{
		ID: uuid.New(),
		Name: param.Name,
		Email: param.Email,
		PasswordHash: hashPass ,
	})
	if err != nil{
		c.IndentedJSON(http.StatusNotAcceptable, gin.H{"msg":"error creating user"})
		log.Printf("unable to create user: %v \n", err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"msg":"created user successfully"})
}


// RETRIEVING ALL USERS

// func(cfg *apiConfig) getAllUsers(c *gin.Context){
// 	users, err := cfg.DB.GetAllUsers(c)
// 	if err != nil{
// 		c.IndentedJSON(http.StatusInternalServerError, gin.H{"msg":"cannot retrieve users"})
// 		log.Printf("error retrieving users %v \n", err)
// 		return
// 	}

// 	c.IndentedJSON(http.StatusFound, users)
// }
