package handlers

import (
	"drivex/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)


type FileHAndler struct {
	fileService *services.FileService
}

func NewFileHandler( fileService *services.FileService) *FileHAndler{
	return &FileHAndler{
		fileService: fileService,
	}
}

type UploadURLRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

func (h * FileHAndler) GenerateUpladUrl(c *gin.Context){
	var req UploadURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file_name and content_type are required",
		})
		return
	}

	userId , exist:= c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}


	uploadUrl , storageKey , err := h.fileService.GenerateUploadURL(c.Request.Context(),
			userId.(string),
			req.FileName,
			req.ContentType,)

			if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate upload URL",
		})
		return
	}

		c.JSON(http.StatusOK, gin.H{
		"upload_url":  uploadUrl,
		"storage_key": storageKey,
	})

}