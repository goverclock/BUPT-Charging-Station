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
	server.POST("/register/user", register_user)
	server.POST("/charge/submit", authMiddleware, charge_submit)
	server.POST("/charge/getChargingMsg", authMiddleware, charge_getChargingMsg)
	// server.POST("/charge/chargeSubmit")	// maybe changeSubmit
	// server.POST("/charge/cancelCharge")
	server.POST("/charge/startCharge", authMiddleware, charge_startCharge)
	// server.POST("/charge/endCharge")
	server.POST("/charge/details", authMiddleware, charge_details)
	server.POST("/recharge", authMiddleware, recharge)
	server.POST("/getbalance", authMiddleware, getbalance)
	// server.POST("/chargeports/getreport")
	// server.POST("/chargeports/getchargeports")
	// server.POST("/chargeports/addchargeport")
	// server.POST("/chargeports/delBatch")
	server.POST("/chargeports/turnon", authMiddleware, chargeports_turnon)
	// server.POST("/chargeports/setfailure")
	// server.POST("/chargeports/waitingCars")
	// server.POST("/system/getsettings")
	// server.POST("/system/setsettings")
	
	server.Run(":8080")
}
