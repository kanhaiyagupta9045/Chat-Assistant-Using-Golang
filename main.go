package main

import (
	"log"
	"os"
	_ "test/config"
	"test/databases"
	"test/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	err := databases.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	routes.FileRoutes(router)

	router.Run(os.Getenv("PORT"))

}
