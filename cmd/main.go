package main

import (
	"log"
	"os"

	// "github.com/ErebusAJ/expense-manager/config"
	"github.com/ErebusAJ/expense-manager/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	godotenv.Load()
	r := gin.Default()

	handlers.RegisterRoutes(r)
	port := os.Getenv("PORT_NO")

	log.Println("Starting server...")

	r.Run("localhost:"+port)

}