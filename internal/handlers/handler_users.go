package handlers

import (
	"log"
	"net/http"

	"github.com/ErebusAJ/expense-manager/internal/db"
	"github.com/ErebusAJ/expense-manager/internal/utils"
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
		c.IndentedJSON(http.StatusNotAcceptable, utils.MessageObj("error creating user"))
		log.Printf("unable to create user: %v \n", err)
		return
	}

	c.IndentedJSON(http.StatusCreated, utils.MessageObj("created user successfully"))
}


// RETRIEVING AUTHENTICATED USERS DETAILS
func(cfg *apiConfig) getAuthUser(c *gin.Context){
	userId, exists := c.Get("userID")
	if(!exists){
		c.JSON(http.StatusUnauthorized, utils.MessageObj("unaauthorized"))
		log.Println("unable to authorize user")
		return
	}

	user, err := cfg.DB.GetUserByID(c, userId.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.MessageObj("error fetching user"))
		log.Printf("unable to fetch user data %v \n", err)
		return 
	}

	c.IndentedJSON(http.StatusOK, user)
}

// LOGIN 
func(cfg *apiConfig) loginUser(c *gin.Context){
	var loginInput struct{
		Email		string		`json:"email" binding:"required,email"`
		Password	string		`json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&loginInput)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.MessageObj("error parsing json data"))
		log.Printf("error parsing json data %v \n", err)
		return 
	}

	email := loginInput.Email
	if email == ""{
		c.IndentedJSON(http.StatusBadRequest, utils.MessageObj("email error"))
		log.Println("error parsing email form body")
		return
	}

	password := loginInput.Password
	if password == ""{
		c.IndentedJSON(http.StatusBadRequest, utils.MessageObj("password error"))
		log.Println("error parsing password from email")
		return
	}

	var dbDetails db.GetUserByEmailRow
	dbDetails, err = cfg.DB.GetUserByEmail(c, email)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.MessageObj("error retrieving password"))
		log.Printf("error retrieving password form DB %v \n", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbDetails.PasswordHash), []byte(password))
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, utils.MessageObj("failed to validate user"))
		log.Printf("failed to validate user %v \n", err)
		return 
	}

	token, err := utils.GenerateJWT(dbDetails.ID)
	if err != nil{
		c.IndentedJSON(http.StatusInternalServerError, utils.MessageObj("error generating token"))
		log.Printf("error generating token %v \n", err)
		return
	}

	c.IndentedJSON(http.StatusOK, token)

}

// RETRIEVING ALL USERS

// func(cfg *apiConfig) getAllUsers(c *gin.Context){
// 	users, err := cfg.DB.GetAllUsers(c)
// 	if err != nil{
// 		c.IndentedJSON(http.StatusInternalServerError, utils.MessageObj("cannot retrieve users"))
// 		log.Printf("error retrieving users %v \n", err)
// 		return
// 	}

// 	c.IndentedJSON(http.StatusFound, users)
// }
