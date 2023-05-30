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
	server.POST("/charge/submit", auth_middleware, charge_submit)
	server.POST("/charge/getChargingMsg", auth_middleware, charge_getChargingMsg)
	server.POST("/charge/chargeSubmit", auth_middleware, charge_changeSubmit)	// maybe changeSubmit
	server.POST("/charge/changeSubmit", auth_middleware, charge_changeSubmit)
	server.POST("/charge/cancelCharge", auth_middleware, charge_cancelCharge)
	server.POST("/charge/startCharge", auth_middleware, charge_startCharge)
	server.POST("/charge/endCharge", auth_middleware, charge_end_charge)
	server.POST("/charge/details", auth_middleware, charge_details)
	server.POST("/recharge", auth_middleware, recharge)
	server.POST("/getbalance", auth_middleware, getbalance)
	// server.POST("/chargeports/getreport")
	// server.POST("/chargeports/getchargeports")
	// server.POST("/chargeports/addchargeport")
	// server.POST("/chargeports/delBatch")
	server.POST("/chargeports/turnon", auth_middleware, chargeports_turnon)
	// server.POST("/chargeports/setfailure")
	// server.POST("/chargeports/waitingCars")
	// server.POST("/system/getsettings")
	// server.POST("/system/setsettings")
	
	server.Run(":8080")
}
