package handlers

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/ErebusAJ/expense-manager/internal/db"
	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// addGroupExpense
// Adds an expense to a specific group and splits among group expense members
// Expects a json body of expense details and a list of participants included in payment
// Uses two sql query functions AddGroupExpense and AddGroupExpenseMembers
func(cfg *apiConfig) addGroupExpense(c *gin.Context){
	var reqDetails struct{
		Title		string	`json:"title" binding:"required"`
		Description	string	`json:"description"`
		Amount		string	`json:"amount" binding:"required"`
		Participants []struct{
			ID	uuid.UUID	`json:"userID" binding:"required"`
			Amount	string	`json:"amount" binding:"required"`
		} `json:"participants" binding:"required"`
	}
	// Bind JSON
	err := c.ShouldBindJSON(&reqDetails)
	if err != nil{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, err)
		return 
	}

	// Check if the split amount among members equals to 
	var sum float64
	sum = 0
	for _, item := range reqDetails.Participants{
		val, _ := strconv.ParseFloat(item.Amount, 64)
		sum = sum + val
	}
	val, _ := strconv.ParseFloat(reqDetails.Amount, 64)
	if(sum != val){
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.AmountSplitError, nil)
		return 
	}

	// parsing userID and groupID
	tempUID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 500, utils.InternalError, utils.MiddlewareError, nil)
		return 
	}
	userID := tempUID.(uuid.UUID)

	tempGID := c.Param("group_id")
	groupID := uuid.MustParse(tempGID)

	// description sql.NullString
	descriptionVal := sql.NullString{
		String: reqDetails.Description,
		Valid: reqDetails.Description != "",
	}

	expense, err := cfg.DB.AddGroupExpense(c, db.AddGroupExpenseParams{
		Title: reqDetails.Title,
		Description: descriptionVal,
		GroupID: groupID,
		CreatedBy: userID,
		Amount: reqDetails.Amount,
	})
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return 
	}

	// Add members recurrently based on request participants
	for _ , member := range reqDetails.Participants{
		_, err := cfg.DB.AddGroupExpenseMembers(c, db.AddGroupExpenseMembersParams{
			GroupExpenseID: expense.ID,
			UserID: member.ID,
			Amount: member.Amount,
		})
		if err != nil{
			utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
			return
		}
	}

	// add members debt or net balance owes/owed
	for _, member := range reqDetails.Participants{
		_, err := cfg.DB.UpdateUserDebts(c, db.UpdateUserDebtsParams{
			FromUser: member.ID,
			ToUser: userID,
			GroupID: groupID,
			Amount: member.Amount,
			ExpenseID: expense.ID,
		})
		if err != nil{
			utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
			return
		}
	} 

	c.IndentedJSON(200, utils.MessageObj("successfully added"))
}


// getAllGroupExpenses
// Retrieves all the expense of a specific group
func(cfg *apiConfig) getAllGroupExpenses(c *gin.Context){
	tempGID := c.Param("group_id")
	groupID := uuid.MustParse(tempGID)

	expenses, err := cfg.DB.GetAllGroupExpenses(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(201, expenses)
}


// updateGroupExpense
// updates a specific group expense 
// uses sql query functions GetGroupExpenseByID, UpdateGroupExpense, GetGroupExpenseMembers
func(cfg *apiConfig) updateGroupExpense(c *gin.Context){
	var reqDetails struct{
		Title		string	`json:"title"`
		Description	string	`json:"description"`
		Amount		string	`json:"amount"`	
	}

	err := c.ShouldBindJSON(&reqDetails)
	if err != nil{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, err)
		return 
	} 

	tempEID := c.Param("expense_id")
	expID := uuid.MustParse(tempEID)
	
	// Retrive expense details
	expense, err := cfg.DB.GetGroupExpenseByID(c, expID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, nil)
		return
	}

	// Get userID of the request sender
	tempUID, exsits := c.Get("userID")
	if !exsits {
		utils.ErrorJSON(c, 500, utils.InvalidError, utils.MiddlewareError, nil)
		return
	}
	userID := tempUID.(uuid.UUID)

	// check if request sender is creator only then update 
	if(userID != expense.CreatedBy){
		utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.InvalidAcess, nil)
		return
	} 

	// Check if fields empty
	amountCheck := true
	if reqDetails.Title == ""{
		reqDetails.Title = expense.Title
	}
	if reqDetails.Description == ""{
		reqDetails.Description = expense.Description.String
	}
	if reqDetails.Amount == ""{
		reqDetails.Amount = expense.Amount
		amountCheck = false
	}

	descriptionVal := sql.NullString{
		String: reqDetails.Description,
		Valid: reqDetails.Description != "",
	}
	// Update Details in DB
	err = cfg.DB.UpdateGroupExpense(c, db.UpdateGroupExpenseParams{
		Title: reqDetails.Title,
		Description: descriptionVal,
		Amount: reqDetails.Amount,
		ID: expID,
	})
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return 
	}

	if amountCheck {
		amount, _ := strconv.ParseFloat(reqDetails.Amount, 64)

		members, err := cfg.DB.GetGroupExpenseMembersByID(c, expID)
		if err != nil{
			utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
			return 
		}
		
		number := len(members)
		equalAmount := amount / float64(number)
		finalAmount := strconv.FormatFloat(equalAmount, 'f', 2, 64)
		for _, member := range members{
			err := cfg.DB.UpdateGroupExpenseMembers(c, db.UpdateGroupExpenseMembersParams{
				Amount: finalAmount,
				ID: member.ID,
			})
			if err != nil{
				utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
				return 
			}
		}
	}

	c.IndentedJSON(204, utils.MessageObj("successfully updated"))
}

// deleteGroupExpense
// deletes a group expense from db by specific id
func(cfg *apiConfig) deleteGroupExpense(c *gin.Context){
	tempEID := c.Param("expense_id")
	expID := uuid.MustParse(tempEID)

	tempUID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.MiddlewareError, nil)
		return
	}
	userID := tempUID.(uuid.UUID)

	expense, err := cfg.DB.GetGroupExpenseByID(c, expID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	if(expense.CreatedBy != userID){
		utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.InvalidAcess, nil)
		return
	}

	err = cfg.DB.DeleteGroupExpense(c, expID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(204, utils.MessageObj("deleted"))

}	


// fetchNetBalance
// retrieves the net balance of a person if he is in +ve or -ve
func(cfg *apiConfig) fetchNetBalance(c *gin.Context){
	tempGID := c.Param("group_id")
	groupID := uuid.MustParse(tempGID)

	netBalances, err := cfg.DB.FetchNetBalance(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return 
	}
	
	c.IndentedJSON(200, netBalances)
}


// minimizeTransactions
// return minimum transactions to make for settling up debts
func(cfg *apiConfig) minimizeTransactions(c *gin.Context){
	tempGID := c.Param("group_id")
	groupID := uuid.MustParse(tempGID)

	netBalances, err := cfg.DB.FetchNetBalance(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	newBalances := make(map[uuid.UUID]string)

	for _, item  := range netBalances{
		newBalances[item.UserID] = item.Netbalance
	}

	transactions := utils.MinimizeDebts(newBalances)

	for _, record := range transactions{
		_, err := cfg.DB.AddSimplifiedTransaction(c, db.AddSimplifiedTransactionParams{
			GroupID: groupID,
			FromUser: record.FromUserID,
			ToUser: record.ToUserID,
			Amount: record.Amount,
		})
		if err != nil{
			utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
			return
		}
	}

	c.IndentedJSON(201, utils.MessageObj("successfully minimized"))
}


// fetchMinimizedTransactions
// Retrieve the minimized transactions from DB
func(cfg *apiConfig) fetchMinimizedTransactions(c *gin.Context){
	tempGID := c.Param("group_id")
	groupID := uuid.MustParse(tempGID)

	records, err := cfg.DB.GetSimplifiedTransactions(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return 
	}

	c.IndentedJSON(200, records)
}

// getGroupTotalExpense
// retrieves a groups total expense
func(cfg *apiConfig) getGroupTotalExpense(c *gin.Context){
	tempGID := c.Param("group_id")
	groupID, err := uuid.Parse(tempGID)
	if err != nil{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, err)
		return
	}

	data, err := cfg.DB.GetTotalGroupExpense(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return
	}

	c.IndentedJSON(200, gin.H{"total_expense":data})
}


// getGroupMembersTotal
// retrieves each members total contribution for the expense group specified
func(cfg *apiConfig) getGroupMembersTotal(c *gin.Context){
	tempUID, exists := c.Get("userID")
	if !exists {
		utils.ErrorJSON(c, 500, utils.InternalError, utils.MiddlewareError, nil)
		return 
	}
	userID := tempUID.(uuid.UUID)


	tempGID := c.Param("group_id")
	groupID, err := uuid.Parse(tempGID)
	if err != nil{
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, err)
		return
	}

	_, err =cfg.DB.CheckMemeber(c, db.CheckMemeberParams{
		GroupID: groupID,
		UserID: userID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows"){
			utils.ErrorJSON(c, 401, utils.UnauthorizedError, utils.InvalidAcess, err)
			return
		}else{
			utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
			return 
		}
	}
	

	data, err := cfg.DB.GetMembersTotalExpense(c, groupID)
	if err != nil{
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return 
	}

	c.IndentedJSON(200, data)
}

// getGroupExpenseDetails
// fetchs details(participants, shares) of a groupExpense specified via ID
func(cfg *apiConfig) getGroupExpenseDetails(c *gin.Context){
	// tempGID := c.Param("group_id")
	// groupID, err1 := uuid.Parse(tempGID)

	tempEID := c.Param("expense_id")
	expenseID, err := uuid.Parse(tempEID)

	if err != nil {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.IDParseError, err)
		return 
	}

	data, err := cfg.DB.GetGroupExpenseDetails(c, expenseID)
	if err != nil {
		utils.ErrorJSON(c, 500, utils.InternalError, utils.DatabaseError, err)
		return 
	}

	c.IndentedJSON(200, data);

}