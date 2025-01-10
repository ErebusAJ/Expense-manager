package handlers

import (
	"log"

	"github.com/ErebusAJ/expense-manager/internal/db"
	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type apiConfig struct{
	DB *db.Queries
}


//User Routes
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

	//Routes
	// r.GET("/auth/user", apiCfg.getAllUsers)
	r.POST("/auth/user", apiCfg.addUser)
}
