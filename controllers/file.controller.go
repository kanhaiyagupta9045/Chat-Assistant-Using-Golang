package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	_ "test/config"
	"test/databases"
	"test/models"

	"github.com/gin-gonic/gin"
)

func FileUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}
		uploadDir := "uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory"})
			return
		}
		filePath := filepath.Join(uploadDir, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		file.Open()
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}
		db := databases.GetDB()
		fileData := models.File{
			Filename: file.Filename,
			Content:  fileContent,
		}

		if err := db.Create(&fileData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": file.Filename})

	}
}
