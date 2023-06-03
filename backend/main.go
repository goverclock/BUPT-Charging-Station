package main

import (
	"buptcs/data"
	_ "net/http/pprof"
	"strconv"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var amazing_lock sync.Mutex

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
	server.POST("/charge/getChargingMsg", auth_middleware, charge_getChargingMsg)	// pending
	server.POST("/charge/chargeSubmit", auth_middleware, charge_changeSubmit)	// maybe changeSubmit	// pending
	server.POST("/charge/changeSubmit", auth_middleware, charge_changeSubmit)
	server.POST("/charge/cancelCharge", auth_middleware, charge_cancelCharge)
	server.POST("/charge/startCharge", auth_middleware, charge_startCharge)
	server.POST("/charge/endCharge", auth_middleware, charge_endCharge)
	server.POST("/charge/details", auth_middleware, charge_details)
	server.POST("/recharge", auth_middleware, recharge)
	server.POST("/getbalance", auth_middleware, getbalance)
	server.POST("/chargeports/getreport", auth_middleware, chargeports_getreport)
	server.POST("/chargeports/getreports", auth_middleware, chargeports_getreports)
	// server.POST("/chargeports/addchargeport")
	// server.POST("/chargeports/delBatch")
	server.POST("/chargeports/switch", auth_middleware, chargeports_switch)
	server.POST("/chargeports/switchBroken", auth_middleware, chargeports_switchBroken)
	server.POST("/chargeports/waitingCars", auth_middleware, chargeports_waitingCars)
	server.POST("/system/getsettings", auth_middleware, system_getsettings)
	server.POST("/system/setsettings", auth_middleware, system_setsettings)
	
	server.Run(":" + strconv.Itoa(data.Port))
}
