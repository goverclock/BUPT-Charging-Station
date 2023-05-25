package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))

	server.POST("/login/user", login_user)
	server.POST("/login/admin", login_admin)
	server.POST("/register/user", register_user)
	server.POST("/charge/submit", authMiddleware, charge_submit)
	server.GET("/charge/getChargingMsg", authMiddleware, charge_getChargingMsg)

	server.Run(":8080")
}
