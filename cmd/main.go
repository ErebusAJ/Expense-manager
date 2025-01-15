package main

import (
	"log"
	"os"

	"github.com/ErebusAJ/expense-manager/internal/handlers"
	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	//Load env
	godotenv.Load()

	//Connect to DB
	DB, err := utils.ConnectDB()
	if err != nil{
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer DB.Close()

	//Initialize router
	r := gin.Default()

	//Register handlers routes
	handlers.RegisterRoutes(r)
	port := os.Getenv("PORT_NO")

	log.Printf("\t\t Starting server... \n")

	r.Run("localhost:"+port)

}