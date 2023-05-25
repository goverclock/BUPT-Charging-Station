package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// almost same as main.go
func init() {
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
	server.GET("/charge/getChargingMsg", authMiddleware, charge_getChargingMsg)
	// server.POST("/charge/chargeSubmit")	// maybe changeSubmit
	// server.POST("/charge/cancelCharge")
	// server.POST("/charge/startCharge")
	// server.POST("/charge/endCharge")
	// server.POST("/charge/details")
	// server.POST("/recharge")
	// server.POST("//getbalance")
	// server.POST("/chargeports/getreport")
	// server.POST("/chargeports/getchargeports")
	// server.POST("/chargeports/addchargeport")
	// server.POST("/chargeports/delBatch")
	// server.POST("/chargeports/turnon")
	// server.POST("/chargeports/setfailure")
	// server.POST("/chargeports/waitingCars")
	// server.POST("/system/getsettings")
	// server.POST("/system/setsettings")

	go server.Run(":8080")
}

func TestLoginUser(t *testing.T) {

}


// assuming user "w", "w" exists
func TestRegisterUser(t *testing.T) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var response struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Username        string  `json:"username"`
		} `json:"data"`
	}
	request.Username = "w"
	request.Password = "w"
	req_body, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", "http://localhost:8080/login/user", bytes.NewBuffer(req_body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &response)
	t.Log(response)
}
