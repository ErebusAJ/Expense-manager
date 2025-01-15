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
	r.POST("/v1/user", apiCfg.addUser)
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
		protected.DELETE("/user", apiCfg.deleteUser)
	} 

	//ADMIN PROTECTED ROUTES
	r.POST("/admin/login", apiCfg.adminLogin)
	adminProtected := r.Group("/admin")
	adminProtected.Use(middleware.AdminMiddleware(signKey))
	{
		adminProtected.GET("/users", apiCfg.adminGetAllUsers)
	}
	// r.GET("/auth/user",)
}
