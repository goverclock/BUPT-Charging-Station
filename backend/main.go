package main

import (
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

	server.POST("/user/login", user_login)
	server.POST("/register/user", register_user)
	server.POST("/charge/submit", authMiddleware, charge_submit)

	server.Run(":8080")
}
