package handlers

import (
	"net/http"
	"os"

	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Handlers for admin

// ADMIN LOGIN
// Calls userLogin function
func(cfg *apiConfig) adminLogin(c *gin.Context){
	var InputDetails struct{
		Email		string	`json:"email" binding:"required,email"`
		Password	string	`json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&InputDetails)
	if err != nil{
		utils.ErrorJSON(c, http.StatusBadRequest, "invalid json body", "error binding request data", err)
		return
	}

	email := InputDetails.Email
	password := InputDetails.Password

	user, err := cfg.DB.GetUserByEmail(c, email)
	if err != nil{
		utils.ErrorJSON(c, http.StatusBadRequest, "invalid user", "error retrieving user from db", err)
		return 
	}

	temp_adminID := os.Getenv("ADMIN_ID")
	adminID, err := uuid.Parse(temp_adminID)
	if err != nil{
		utils.ErrorJSON(c, http.StatusInternalServerError, "internal error", "error casting adminID", err)
		return
	}

	adminRole := os.Getenv("ADMIN_ACCESS_LEVEL")
	if(user.AccessLevel != adminRole || user.ID != adminID){
		utils.ErrorJSON(c, http.StatusUnauthorized, "unauthorized", "inavlid access level", nil)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil{
		utils.ErrorJSON(c, http.StatusBadRequest, "failed to validate user", "inavlid user password", err)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.AccessLevel)
	if err != nil{
		utils.ErrorJSON(c, http.StatusInternalServerError, "unable to validate login", "error generating token", err)
		return
	}

	c.IndentedJSON(http.StatusOK, token)
	
}

// GET ALL USER DETAILS
func(cfg *apiConfig) adminGetAllUsers(c *gin.Context){
	// Retreive userId, userRole from middleware
	userID, existsI := c.Get("userId")
	userRole, existsR := c.Get("userRole")
	if !existsI || !existsR{
		utils.ErrorJSON(c, http.StatusUnauthorized, "unauthorized", "unable to retrieve user/role from middleware", nil)
		return
	}

	godotenv.Load()
	adminID := os.Getenv("ADMIN_ID")
	adminRole := os.Getenv("ADMIN_ACCESS_LEVEL")
	if adminID == "" || adminRole == ""{
		utils.ErrorJSON(c, http.StatusInternalServerError, "internal server error", "error retrieving adminID/adminRole", nil)
		return
	}

	if userID != adminID || userRole != adminRole{
		utils.ErrorJSON(c, http.StatusUnauthorized, "unauthorized access", "error invalid access level", nil)
		return
	}

	users, err := cfg.DB.GetAllUsers(c)
	if err != nil{
		utils.ErrorJSON(c, http.StatusInternalServerError, "error fetching users", "error fetching users details", err)
		return 
	}

	c.IndentedJSON(http.StatusOK, users)
}	
