package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length"},
	}))

	server.POST("/user/login", func(ctx *gin.Context) {
		var js struct {
			Username string	`json:"username"`
			Password string	`json:"password"`
		}
		ctx.BindJSON(&js)
		log.Println(js)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	server.Run(":8080")
}
