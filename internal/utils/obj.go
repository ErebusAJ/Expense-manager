package utils

import "github.com/gin-gonic/gin"

func MessageObj(msg string)(map[string]any){
	return gin.H{"message":msg}
}