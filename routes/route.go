package routes

import (
	"test/controllers"

	"github.com/gin-gonic/gin"
)

func FileRoutes(router *gin.Engine) {
	router.POST("/upload",controllers.FileUpload())
}
