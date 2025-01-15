package utils

import(
	"log"
	
	"github.com/gin-gonic/gin"
)

// CUSTOM ERROR
// This is function when we need to return err to client from a handler
// Takes a gin.Context, two strings(client, server), error
func ErrorJSON(c *gin.Context,code int, client string, server string, err error){
	c.IndentedJSON(code, MessageObj(client))
	log.Printf(server + " %v", err)
}