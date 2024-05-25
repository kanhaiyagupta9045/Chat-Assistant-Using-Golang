package controllers

// 9155340847
import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"test/chatassistant"
	_ "test/config"
	"test/databases"
	"test/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func FileUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request models.FileUploadRequest
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}
		assistantData := c.PostForm("assistant")
		if err := json.Unmarshal([]byte(assistantData), &request.Assistant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assistant data"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(request.Assistant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open the file"})
			return
		}
		defer src.Close()

		filecontent, err := ioutil.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read the file"})
			return
		}
		filedata := models.File{
			Filename: file.Filename,
			Content:  filecontent,
		}

		db := databases.GetDB()
		if db == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to database"})
			return
		}
		if err := db.Create(&filedata).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erorr": err.Error()})
			return
		}

		msg, err := chatassistant.ChatAssistant(file.Filename, filecontent, string(request.Assistant.AssistantName), string(request.Assistant.AssistantInstructions), string(request.Assistant.Content))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": msg, "database_msg": "File Uploaded Succssfully"})
	}
}
