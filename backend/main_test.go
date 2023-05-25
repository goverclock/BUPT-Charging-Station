package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var authToken string

// almost same as main.go
func init() {
	return
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

	go server.Run(":8080")	// only difference
}

// no Authorization header
func send(method string, route string, request interface{}) []byte {
	route = "http://localhost:8080" + route
	req_body, _ := json.Marshal(request)
	req, err := http.NewRequest(method, route, bytes.NewBuffer(req_body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if token := resp.Header.Get("Authorization"); token != "" {
		authToken = token
		fmt.Println("Token is", authToken)
	}
	
	body, _ := io.ReadAll(resp.Body)
	return body
}

// assuming user "w", "w" exists
func TestLoginUser(t *testing.T) {
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

	body := send("POST", "/login/user", request)
	json.Unmarshal(body, &response)
	t.Log(response)
	if response.Code != 200 {
		t.Fail()
	}
}

// func TestRegisterUser(t *testing.T) {
// }

// /charge/submit
func TestChargeSubmit(t *testing.T) {
	TestLoginUser(t)	// login in as user "w"
	
	var request struct {
		ChargeMode   int     `json:"chargeMode"`
		ChargeAmount float64 `json:"chargeAmount"`
	}
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	request.ChargeMode = 0
	request.ChargeAmount = 2.0

	body := send("POST", "/charge/submit", request)
	json.Unmarshal(body, &response)
	t.Log(response)
	if response.Code != 200 {
		t.Fail()
	}
}
