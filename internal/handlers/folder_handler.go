package handlers

import (
	"drivex/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FolderHandler struct {
	folderService *services.FolderService
}

func NewFolderHandler(
	folderService *services.FolderService,
) *FolderHandler {
	return &FolderHandler{
		folderService: folderService,
	}
}
type UpdateFolderRequest struct {
	Name string `json:"name" binding:"required"`
}


type CreateFolderRequest struct {
	Name     string  `json:"name" binding:"required"`
	ParentID *string `json:"parent_id"`
}

func (h *FolderHandler) Create(c *gin.Context) {

	var req CreateFolderRequest

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

	folder, err := h.folderService.Create(
		c.Request.Context(),
		userID.(string),
		req.ParentID,
		req.Name,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, folder)
}



func ( h *FolderHandler) GetFolders( c *gin.Context){
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var parentId *string

		parentIDQuery := c.Query("parent_id")

	if parentIDQuery != "" {
		parentId = &parentIDQuery
	}

	folders , err := h.folderService.GetFolders(
		c.Request.Context(),
		userID.(string),
		parentId,
	)
		if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch folders",
		})
		return
	}

		c.JSON(http.StatusOK, gin.H{
		"folders": folders,
	})
}




func (h *FolderHandler) Update(c *gin.Context) {

	var req UpdateFolderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "folder name is required",
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

	folder, err := h.folderService.Update(
		c.Request.Context(),
		userID.(string),
		c.Param("id"),
		req.Name,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, folder)
}


func (h *FolderHandler) Delete(c *gin.Context) {

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	err := h.folderService.Delete(
		c.Request.Context(),
		userID.(string),
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete folder",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "folder deleted successfully",
	})
}