package controllers

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"test/chatassistant"
	"test/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ChatController() gin.HandlerFunc {
	return func(c *gin.Context) {
		var assistant models.Assistant

		if err := c.BindJSON(&assistant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validate.Struct(assistant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		uploadsDir := "uploads"
		if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{"message": "No Upload directory found"})
			return
		}
		files, err := ioutil.ReadDir(uploadsDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read uploads directory"})
			return
		}
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please Upload the files no files found"})
			return
		}

		for _, file := range files {
			if !file.IsDir() {
				filePath := filepath.Join(uploadsDir, file.Name())
				fileContent, err := ioutil.ReadFile(filePath)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
					return
				}
				msg, err := chatassistant.ChatAssistant(file.Name(), fileContent, assistant.AssistantName, assistant.AssistantInstructions, assistant.Content)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"erorr": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"msg": msg})
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error occured While chatting with chat Assistant"})

	}
}
