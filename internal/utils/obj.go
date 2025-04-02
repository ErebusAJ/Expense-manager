package utils

import "github.com/gin-gonic/gin"

// Client Errors
const (
	InternalError		=	"Internal Error"
	InvalidError		=	"Invalid Request"
	UnauthorizedError	=	"Unauthorized"
	NotFound			=	"Not Found"
)

// Server Log Errors
const (
	RequestBodyError	=	"error binding request json"
	MiddlewareError		=	"error retreiving user id from middleware"
	DatabaseError		= 	"error processing db query"
	IDParseError		=	"error parsing to uuid"
	InvalidAcess		=	"error not appropriate access level"
	AmountSplitError	=	"error wrong amount split"

)


func MessageObj(msg string)(map[string]any){
	return gin.H{"message":msg}
}