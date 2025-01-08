package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine){
	r.GET("auth/get-albums", getAlbums)
}
