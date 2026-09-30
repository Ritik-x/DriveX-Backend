package handlers

import (
	"drivex/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)


type FileShareHandler struct {
	service *services.FileShareService
}

func NewFileShareHandler(
	service *services.FileShareService ,
	
 ) *FileShareHandler{

	return &FileShareHandler{
		service: service,
	}
 }

 type ShareFileRequest struct {

	Email      string `json:"email" binding:"required,email"`
	Permission string `json:"permission" binding:"required"`
 }


 func (h *FileShareHandler) ShareFile(c *gin.Context) {

	userID := c.GetString("user_id")
	fileID := c.Param("id")

	var req ShareFileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.ShareFile(
		c.Request.Context(),
		userID,
		fileID,
		req.Email,
		req.Permission,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "file shared successfully",
	})
}