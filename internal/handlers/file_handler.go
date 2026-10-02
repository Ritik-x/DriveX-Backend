package handlers

import (
	"drivex/internal/models"
	"drivex/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileHAndler struct {
	fileService   *services.FileService
	folderService *services.FolderService
}

func NewFileHandler(fileService *services.FileService, folderService *services.FolderService) *FileHAndler {
	return &FileHAndler{
		fileService:   fileService,
		folderService: folderService,
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

func (h *FileHAndler) GenerateUpladUrl(c *gin.Context) {
	var req UploadURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file_name and content_type are required",
		})
		return
	}

	userId, exist := c.Get("user_id")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	uploadUrl, storageKey, fileId, err := h.fileService.GenerateUploadURL(
		c.Request.Context(),
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
		"file_id":     fileId,
		"upload_url":  uploadUrl,
		"storage_key": storageKey,
	})

}

func (h *FileHAndler) CompleteUpload(c *gin.Context) {
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

	file, err := h.fileService.CompleteUpload(
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
func (h *FileHAndler) GetFiles(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var folderId *string
	queryFolderID := c.Query("folder_id")

	if queryFolderID != "" {
		folderId = &queryFolderID
	}

	files, err := h.fileService.GetFiles(
		c.Request.Context(),
		userId.(string),
		folderId,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch files",
		})
		return
	}

	if files == nil {
		files = []models.File{}
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

func (h *FileHAndler) DeleteFile(c *gin.Context) {
	userId := c.GetString("user_id")

	fileId := c.Param("id")

	err := h.fileService.DeletFile(c.Request.Context(), fileId, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file moved to trash",
	})
}

func (h *FileHAndler) GetTrash(c *gin.Context) {
	userId := c.GetString("user_id")

	files, err := h.fileService.GetTrash(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get trash",
		})
		return
	}

	if files == nil {
		files = []models.File{}
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

func (h *FileHAndler) RestoreFile(c *gin.Context) {
	userID := c.GetString("user_id")
	fileID := c.Param("id")

	err := h.fileService.RestoreFile(
		c.Request.Context(),
		fileID,
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to restore file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file restored successfully",
	})
}

func (h *FileHAndler) PermanentDeleteFile(c *gin.Context) {
	userID := c.GetString("user_id")
	fileID := c.Param("id")

	err := h.fileService.PermanentDelete(
		c.Request.Context(),
		fileID,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to permanently delete file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "file permanently deleted",
	})
}

func (h *FileHAndler) Search(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusOK, gin.H{
			"files":   []models.File{},
			"folders": []models.Folder{},
		})
		return
	}

	if len(query) > 100 {
		query = query[:100]
	}

	files, err := h.fileService.Search(c.Request.Context(), userID.(string), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to search files",
		})
		return
	}
	if files == nil {
		files = []models.File{}
	}

	folders, err := h.folderService.Search(c.Request.Context(), userID.(string), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to search folders",
		})
		return
	}
	if folders == nil {
		folders = []models.Folder{}
	}

	c.JSON(http.StatusOK, gin.H{
		"files":   files,
		"folders": folders,
	})
}

func (h *FileHAndler) GetSharedFie(c *gin.Context) {
	userId := c.GetString("user_id")
	files, err := h.fileService.GetSharedFiles(c.Request.Context(), userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get shared files",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}
