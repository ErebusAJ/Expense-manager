package handlers

import (
	"log"
	"os"

	"github.com/ErebusAJ/expense-manager/internal/db"
	"github.com/ErebusAJ/expense-manager/internal/middleware"
	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/joho/godotenv"

	// "github.com/ErebusAJ/expense-manager/internal/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type apiConfig struct{
	DB *db.Queries
}

// ROUTES
func RegisterRoutes(r *gin.Engine){
	//Establishing DB connections for handlers
	DB, err := utils.ConnectDB()
	if err != nil{
		log.Fatalf("failed to connect to DB: %v", err)
	}

	new_db := db.New(DB)
	apiCfg := apiConfig{
		DB: new_db,
	}

	// Routes
	r.POST("/v1/register", apiCfg.registerUser)
	r.POST("/v1/login", apiCfg.loginUser)
	//JWT AUTHENTICATED ROUTES
	// Secret key
	godotenv.Load()
	signKey := os.Getenv("SECRET_KEY")
	if signKey == ""{
		log.Printf("unable to fetch signed key")
	}

	//AUTH ROUTES
	protected := r.Group("/auth")
	protected.Use(middleware.AuthMiddleware(signKey))
	{
		protected.GET("/user", apiCfg.getAuthUser)
		protected.PUT("/user", apiCfg.updateUserDetails)
		protected.DELETE("/user", apiCfg.deleteUser)

		//EXPENSES ROUTES
		protected.GET("/expenses", apiCfg.getUserExpenses)
		protected.GET("/expense/:id", apiCfg.getUserExpenseByID)
		protected.POST("/expense", apiCfg.addUserExpense)
		protected.PUT("/expense/:id", apiCfg.updateUserExpense)
		protected.DELETE("/expense/:id", apiCfg.deleteExpense)
		protected.GET("/expense/total", apiCfg.getTotalExpense)

		//GROUPS ROUTES
		protected.GET("/group/:group_id", apiCfg.getGroupByID)
		protected.POST("/group", apiCfg.userCreateGroup)
		protected.PUT("/group/:group_id", apiCfg.updateGroupDetails)
		protected.DELETE("/group/:group_id", apiCfg.deleteUserGroup)
		
		protected.POST("/group/:group_id/member/:user_id", apiCfg.addGroupMember)
		protected.GET("/group/:group_id/member", apiCfg.getGroupMembers)
		protected.DELETE("/group/:group_id/member/:user_id", apiCfg.deleteGroupMember)
		protected.DELETE("/group/:group_id/member", apiCfg.deleteLoggedInUser)

		protected.GET("/group/all", apiCfg.getUserAllGroups)
		protected.GET("/group", apiCfg.getUserGroups)

		//GROUPS EXPENSES ROUTES
		protected.POST("/group/:group_id/expense", apiCfg.addGroupExpense)
		protected.GET("/group/:group_id/expense", apiCfg.getAllGroupExpenses)
		protected.PUT("/group/:group_id/expense/:expense_id", apiCfg.updateGroupExpense)
		protected.DELETE("/group/:group_id/expense/:expense_id", apiCfg.deleteGroupExpense)

		// GROUP EXPENSE ANALYTICS
		protected.GET("/group/:group_id/expense-total", apiCfg.getGroupTotalExpense)
		protected.GET("/group/:group_id/member-total", apiCfg.getGroupMembersTotal)
		
		//GROUP MEMBERS NET BALANCE
		protected.GET("/group/:group_id/netbalance", apiCfg.fetchNetBalance)
		protected.POST("/group/:group_id/minimizeTransaction", apiCfg.minimizeTransactions)
		protected.GET("/group/:group_id/minimizeTransaction", apiCfg.fetchMinimizedTransactions)


	} 
	// USER PASS RESET
	r.POST("/v1/user/password-reset", apiCfg.resetPasswordRequest)
	r.POST("/v1/user/password-reset/:token", apiCfg.resetPasswordConfirm)

	//ADMIN PROTECTED ROUTES
	r.POST("/admin/login", apiCfg.adminLogin)
	adminProtected := r.Group("/admin")
	adminProtected.Use(middleware.AdminMiddleware(signKey))
	{
		adminProtected.GET("/users", apiCfg.adminGetAllUsers)
		adminProtected.GET("/user/:id", apiCfg.adminGetUserByID)
		adminProtected.DELETE("/user/:id", apiCfg.adminDeleteUserByID)
	}
}
