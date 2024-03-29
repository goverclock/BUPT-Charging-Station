package main

import (
	"buptcs/data"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

var authToken string

// Authorization header included
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
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Username string `json:"username"`
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

func TestLoginAdmin(t *testing.T) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Username string `json:"username"`
		} `json:"data"`
	}
	request.Username = "q"
	request.Password = "q"

	body := send("POST", "/login/user", request)
	json.Unmarshal(body, &response)
	t.Log(response)
	if response.Code != 200 {
		t.Fail()
	}
	t.Log("logined as admin")
}

// /charge/submit
func TestChargeSubmit(t *testing.T) {
	TestLoginUser(t) // login in as user "w"

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

func TestGetChargingMsg(t *testing.T) {
	TestLoginUser(t)

	var request struct {
		Username string `json:"username"`
	}
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Waiting_count int `json:"waiting_count"`
		} `json:"data"`
	}

	request.Username = "w"
	body := send("POST", "/charge/getChargingMsg", request)
	json.Unmarshal(body, &response)
	t.Log(response)
	if response.Code != 200 {
		t.Fail()
	}
}

func TestStartCharge(t *testing.T) {
	TestLoginUser(t)

	var request struct {
		User_id int `json:"user_id"`
	}
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
		} `json:"data"`
	}

	request.User_id = 1
	body := send("POST", "/charge/startCharge", request)
	json.Unmarshal(body, &response)
	t.Log(response)
	if response.Code != 200 {
		t.Fail()
	}
}

func TestDetails(t *testing.T) {
	TestLoginUser(t)

	var request struct {
		User_id int `json:"user_id"`
	}
	var response struct {
		Code int           `json:"code"`
		Msg  string        `json:"msg"`
		Data []data.Report `json:"data"`
	}

	request.User_id = 1
	body := send("POST", "/charge/details", request)
	json.Unmarshal(body, &response)
	t.Log(response)
	if response.Code != 200 {
		t.Fail()
	}
}

func TestFailure(t *testing.T) {
	TestLoginAdmin(t)

	var request struct {
		Charge_id int `json:"charge_id"`
	}
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
		} `json:"data"`
	}

	request.Charge_id = 0
	body := send("POST", "/chargeports/switchBroken", request)
	json.Unmarshal(body, &response)
	t.Log(response)
	if response.Code != 200 {
		t.Fail()
	}
}
