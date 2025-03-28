package handlers

import (
	"database/sql"

	"github.com/ErebusAJ/expense-manager/internal/db"
	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


// addUserExpense
// Adds a user expense to the DB
// doesn't return anything only a success message
func(cfg *apiConfig) addUserExpense(c *gin.Context){
	var reqDetails struct{
		Amount		string	`json:"amount" binding:"required"`
		Title		string	`json:"title" binding:"required"`
		Description	string	`json:"description"`	
	}

	err := c.ShouldBindJSON(&reqDetails)
	if err != nil{
		utils.ErrorJSON(c, 400, "invalid request", "error binding request body", err)
		return
	}

	tempID, exists := c.Get("userID")
	if !exists{
		utils.ErrorJSON(c, 400, "invalid request", "error retrieving user form middleware", nil)
		return
	}

	userID := tempID.(uuid.UUID)
	descriptionVal := sql.NullString{
		String: reqDetails.Description,
		Valid: reqDetails.Description != "",
	}

	err = cfg.DB.AddExpense(c, db.AddExpenseParams{
		UserID: userID,
		Amount: reqDetails.Amount,
		Title: reqDetails.Title,
		Description: descriptionVal,
	})
	if err != nil{
		utils.ErrorJSON(c, 500, "internal error", "error adding expense to db", err)
		return 
	}

	c.IndentedJSON(201, utils.MessageObj("expense added successfully"))
} 


// getUserExpenses
// Returns all the expenses of user
func(cfg *apiConfig) getUserExpenses(c *gin.Context){
	tempID, exists := c.Get("userID")
	if !exists{
		utils.ErrorJSON(c, 400, "invalid request", "error retrieving user form middleware", nil)
		return
	}

	userID := tempID.(uuid.UUID)

	expenseList, err := cfg.DB.GetAllExpenses(c, userID)
	if err != nil{
		utils.ErrorJSON(c, 500, "internal error", "error retrieving expenses from db", err)
		return 
	}

	c.IndentedJSON(200, expenseList)
}


// getUserExpenseByID
// Returns a specific expense of user via id
func(cfg *apiConfig) getUserExpenseByID(c *gin.Context){
	tempID := c.Param("id")
	expID, err := uuid.Parse(tempID)
	if err != nil{
		utils.ErrorJSON(c, 400, "invalid request", "error parsing id to uuid", err)
		return 
	}

	expense, err := cfg.DB.GetExpenseByID(c, expID)
	if err != nil{
		utils.ErrorJSON(c, 500, "internal error", "error retrieving expense from db", err)
		return 
	}

	c.IndentedJSON(200, expense)
}

// updateUserExpense
// update the fields in user expense specified by id 
// returns updated feild
func(cfg *apiConfig) updateUserExpense(c *gin.Context){
	var reqDetails struct{
		Amount		string	`json:"amount"`
		Title		string	`json:"title"`
		Description	string	`json:"description"`
	}

	err := c.ShouldBindJSON(&reqDetails)
	if err != nil{
		utils.ErrorJSON(c, 400, "invalid request", "error binding req json", err)
		return 
	}

	tempID := c.Param("id")
	
	expID, err := uuid.Parse(tempID)
	if err != nil{
		utils.ErrorJSON(c, 400, "invalid request", "error passing id to uuid", err)
		return
	}

	expense, err := cfg.DB.GetExpenseByID(c, expID)
	if err != nil{
		utils.ErrorJSON(c, 500, "internal error", "error retrieving expense details from db", err)
		return 
	}

	// Check which field updated
	// If request field is empty assing existing value 
	if reqDetails.Amount == ""{
		reqDetails.Amount = expense.Amount 
	}
	if reqDetails.Title == ""{
		reqDetails.Title = expense.Title
	}
	if reqDetails.Description == ""{
		reqDetails.Description = expense.Description.String
	}

	descriptionVal := sql.NullString{
		String: reqDetails.Description,
		Valid: reqDetails.Description != "",
	}

	updatedExp, err := cfg.DB.UpdateExpense(c, db.UpdateExpenseParams{
		Amount: reqDetails.Amount,
		Title: reqDetails.Title,
		Description: descriptionVal,
		ID: expID,
	})
	if err != nil{
		utils.ErrorJSON(c, 500, "internal error", "error updating value in db", err)
		return
	}

	c.IndentedJSON(204, updatedExp)
}


// deleteExpense
// Deletes a expense from database specified by id
func(cfg *apiConfig) deleteExpense(c *gin.Context){
	tempID := c.Param("id")

	expID, err := uuid.Parse(tempID)
	if err != nil{
		utils.ErrorJSON(c, 400, "invalid request", "error parsing id to uuid", err)
		return
	}

	err = cfg.DB.DeleteExpense(c, expID)
	if err != nil{
		utils.ErrorJSON(c, 500, "internal error", "error deleting expense from db", err)
		return 
	}

	c.IndentedJSON(204, utils.MessageObj("deletion successful"))
}


// getTotalExpense
// retrieves a user's total expense
func(cfg* apiConfig) getTotalExpense(c *gin.Context){
	tempID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.MiddlewareError, nil)
		return
	}
	userID := tempID.(uuid.UUID)


	data, err := cfg.DB.TotalExpense(c, userID)
	if err != nil {
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(200, gin.H{"total-expense": data})
}