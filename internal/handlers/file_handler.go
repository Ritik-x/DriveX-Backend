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



type CompleteUploadRequest struct {
	FileID       string  `json:"file_id" binding:"required"`
	FolderID     *string `json:"folder_id"`
	FileName     string  `json:"file_name" binding:"required"`
	OriginalName string  `json:"original_name" binding:"required"`
	StorageKey   string  `json:"storage_key" binding:"required"`
	MimeType     string  `json:"mime_type" binding:"required"`
	Size         int64   `json:"size" binding:"required"`
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
			req.ContentType,
		)

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

func (h *FileHAndler) CompleteUpload( c *gin.Context){
	var req CompleteUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

		userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}


	file , err := h.fileService.CompleteUpload(
		c.Request.Context(),
		userID.(string),
		req.FileID,
		req.FolderID,
		req.FileName,
		req.OriginalName,
		req.StorageKey,
		req.MimeType,
		req.Size,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save file metadata",
		})
		return
	}

	c.JSON(http.StatusCreated, file)

}



func (h *FileHAndler) Download(c *gin.Context) {

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	url, err := h.fileService.GenerateDownloadUrl(
		c.Request.Context(),
		userID.(string),
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "file not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"download_url": url,
	})
}