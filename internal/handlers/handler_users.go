package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ErebusAJ/expense-manager/internal/db"
	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CREATING NEW USER
// hashPassword function is used to hash the password recieved
// addUser creates the uses using sqlc DB queries

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (cfg *apiConfig) registerUser(c *gin.Context) {
	type user struct {
		Name     string    `json:"name" binding:"required"`
		Email    string    `json:"email" binding:"required"`
		Password string    `json:"password" binding:"required"`
	}

	var param user
	err := c.BindJSON(&param)
	if err != nil {
		utils.ErrorJSON(c, 400, utils.RequestBodyError, utils.RequestBodyError, err)
		return
	}

	//Hashing Password
	hashPass, err := hashPassword(param.Password)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, "internal server error", "error hashing password", err)
		return
	}

	err = cfg.DB.CreateUser(c, db.CreateUserParams{
		ID:           uuid.New(),
		Name:         param.Name,
		Email:        param.Email,
		PasswordHash: hashPass,
	})
	if strings.Contains(err.Error(), "duplicate key value") {
		utils.ErrorJSON(c, http.StatusNotAcceptable, "email already in use", "unable to create user", err)
		return
	}else if err != nil {
		utils.ErrorJSON(c, http.StatusNotAcceptable, "unable to create user", "unable to create user", err)
		return
	}

	c.IndentedJSON(http.StatusCreated, utils.MessageObj("created user successfully"))
}

// LOGIN USER
// Verifies the incoming payload containing email, password
// If validated returns a JWT Auth token
func (cfg *apiConfig) loginUser(c *gin.Context) {
	var loginInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBindJSON(&loginInput)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, "invalid json body", "error parsing json data", err)
		return
	}

	email := loginInput.Email
	password := loginInput.Password
	if email == "" || password == "" {
		utils.ErrorJSON(c, http.StatusBadRequest, "invalid json body", "error parsing email from body", nil)
		return
	}

	dbDetails, err := cfg.DB.GetUserByEmail(c, email)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, "invalid user", "error user not found", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbDetails.PasswordHash), []byte(password))
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, "validation failed", "invalud user details", err)
		return
	}

	token, err := utils.GenerateJWT(dbDetails.ID, dbDetails.AccessLevel)
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, "error generating token", "error generating token", err)
		return
	}

	c.IndentedJSON(http.StatusOK, token)

}

// RETRIEVING AUTHENTICATED USERS DETAILS
// Authentication is done via a AuthMiddleware using JWT
// Middleware passes the userid map claim as context to the handlers
func (cfg *apiConfig) getAuthUser(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, http.StatusUnauthorized, "unauthorized", "unable to authorize user", nil)
		return
	}

	user, err := cfg.DB.GetUserByID(c, userId.(uuid.UUID))
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, "error fetching user", "unable to fetch user data", err)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

// UPDATE USER DETAILS
// Updated detials based on JSON body
// If a parameter is missing assigns existing value in db
func (cfg *apiConfig) updateUserDetails(c *gin.Context) {
	var InputDetails struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.ShouldBindJSON(&InputDetails)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, "invalid request body", "error binding request json body", err)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, http.StatusInternalServerError, "internal error", "error retrieving user from midlleware", nil)
		return
	}

	userDetails, err := cfg.DB.GetUserByID(c, userID.(uuid.UUID))
	if err != nil {
		utils.ErrorJSON(c, http.StatusInternalServerError, "internal error", "error getting users details from db", err)
		return
	}

	// Checking if required field value exists in JSON body if not assign existing value from DB
	if InputDetails.Name == "" {
		InputDetails.Name = userDetails.Name
	}

	if InputDetails.Email == "" {
		InputDetails.Email = userDetails.Email
	}

	if InputDetails.Password == "" {
		InputDetails.Password = userDetails.PasswordHash
	} else {
		InputDetails.Password, err = hashPassword(InputDetails.Password)
		if err != nil {
			utils.ErrorJSON(c, http.StatusInternalServerError, "internal error", "error hashing password", err)
			return
		}
	}

	// Updating values

	err = cfg.DB.UpdateUserDetails(c, db.UpdateUserDetailsParams{
		Name:         InputDetails.Name,
		PasswordHash: InputDetails.Password,
		Email:        InputDetails.Email,
		ID:           userID.(uuid.UUID),
	})
	if err != nil {
		utils.ErrorJSON(c, 500, "unable to update user", "error updating user in db", err)
		return
	}

	c.IndentedJSON(http.StatusNoContent, utils.MessageObj("successfully updated user"))
}

// DELETE USER IF AUTHENTICATED
// Deletes user from DB return 204 no content
func (cfg *apiConfig) deleteUser(c *gin.Context) {
	tempID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, http.StatusBadRequest, "unauthorized", "unable to get userId from middleware", nil)
		return
	}

	userID := tempID.(uuid.UUID)

	err := cfg.DB.DeleteUserByID(c, userID)
	if err != nil {
		utils.ErrorJSON(c, http.StatusBadRequest, "error deleting user", "error deleting user", err)
		return
	}

	c.IndentedJSON(http.StatusNoContent, utils.MessageObj("successfully deleted user"))
}

// PASSWORD RESET REQUEST
// Sends a password request
// Takes user email verfies user exists or not
// Sends a email with reset link
func (cfg *apiConfig) resetPasswordRequest(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.ErrorJSON(c, 400, "invalid", "error binding request body to JSON", err)
		return
	}

	user, err := cfg.DB.GetUserByEmail(c, req.Email)
	if err != nil {
		utils.ErrorJSON(c, 404, "user not found", "error finding user in db", err)
		return
	}

	token := uuid.NewString()
	expires := time.Now().Add(1 * time.Hour)

	err = cfg.DB.InsertToken(c, db.InsertTokenParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expires,
	})
	if err != nil {
		utils.ErrorJSON(c, 500, "internal error", "error inserting generated token to db", err)
		return
	}

	url := fmt.Sprintf("http:localhost:8080/users/reset-password/%v", token)

	err = utils.SendEmail(req.Email, fmt.Sprintf("Password rest link %v", url))
	if err != nil {
		utils.ErrorJSON(c, 500, "unable to send reset link", "error sending mail to user", err)
		return
	}

	c.IndentedJSON(200, utils.MessageObj("email sent successfully"))
}

// PASSWOR CONFIRM
// Confirms the token from previous request
// Updates password and deletes token from table
func (cfg *apiConfig) resetPasswordConfirm(c *gin.Context) {
	var req struct {
		NewPassword string `json:"password"`
	}

	reqToken := c.Param("token")

	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.ErrorJSON(c, 400, "invalid request", "error binding request body to json", err)
		return
	}

	token, err := cfg.DB.GetUserToken(context.Background(), reqToken)
	if err != nil {
		utils.ErrorJSON(c, 400, "invalid token", "error unable to find token in db", err)
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		_ = cfg.DB.DeleteToken(c, token.Token)

		utils.ErrorJSON(c, 400, "expired token", "token expired", err)
		return
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil{
		utils.ErrorJSON(c, 500, "internal error", "error hashing password", err)
		return
	}

	err = cfg.DB.SetPassword(context.Background(), db.SetPasswordParams{
		PasswordHash: hashedPassword,
		ID: token.UserID,
	})
	if err != nil{
		utils.ErrorJSON(c, 500, "internal error", "error updating password in db", err)
		return
	}

	err = cfg.DB.DeleteToken(context.Background(), reqToken)
	if err != nil{
		utils.ErrorJSON(c, 500, "internal error", "error deleting token from db", err)
		return
	}

	c.IndentedJSON(200, utils.MessageObj("password reset successfull"))
}
