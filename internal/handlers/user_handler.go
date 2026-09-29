package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type UserHandler struct {}

func NewUserHandler() *UserHandler{
	return &UserHandler{}
}


func ( h *UserHandler) Me ( c *gin.Context){

	userID, exists := c.Get("user_id")


	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

		c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
	})
}