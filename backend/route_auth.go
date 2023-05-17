package main

import (
	"buptcs/data"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func user_login(ctx *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ctx.BindJSON(&request)
	log.Println(request)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Username        string  `json:"username"`
			Password        string  `json:"password"`
			Balance         float64 `json:"balance"`
			BatteryCapacity float64 `json:"batteryCapacity"`
		} `json:"data"`
	}

	// description:
	// String CODE_200 = "200"; //成功
	// String CODE_500 = "500"; //系统错误
	// String CODE_400 = "400"; //参数错误
	// String CODE_401 = "401"; //权限不足 TODO
	// String CODE_600 = "600"; //其它业务异常
	// authenticate
	response.Data.Username = request.Username
	response.Data.Password = request.Password
	user, err := data.UserByName(request.Username)
	if request.Username == "" || request.Password == "" {
		response.Code = "400"
		response.Msg = "need user name or password"
	} else if err != nil {
		response.Code = "500"
		response.Msg = "no such user"
	} else if user.Password != data.Encrypt(request.Password) {
		response.Code = "500"
		response.Msg = "wrong password"
	} else if user.Password == data.Encrypt(request.Password) {
		response.Code = "200"
		response.Msg = "success"
	}
	response.Data.Balance = user.Balance
	response.Data.Balance = user.BatteryCapacity

	ctx.JSON(http.StatusOK, response)
}

func register_user(ctx *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ctx.BindJSON(&request)
	log.Println(request)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Username        string  `json:"username"`
			Password        string  `json:"password"`
			Balance         float64 `json:"balance"`
			BatteryCapacity float64 `json:"batteryCapacity"`
		} `json:"data"`
	}

	user, err := data.UserByName(request.Username)
	if request.Username == "" || request.Password == "" {
		response.Code = "400"
		response.Msg = "need user name or password"
	} else if err == nil {
		// user with same name already exists, can't register
		response.Code = "600"
		response.Msg = "username is taken"
	} else {
		response.Code = "200"
		response.Msg = "register succeeds"
		user.Name = request.Username
		user.Password = request.Password
		user.Balance = 0
		user.BatteryCapacity = 0
		err = user.Create()
		if err != nil {
			log.Fatal(err, "fail to create user")
		}

		response.Data.Username = user.Name
		response.Data.Password = user.Password
		response.Data.Balance = user.Balance
		response.Data.BatteryCapacity = user.BatteryCapacity
	}

	ctx.JSON(http.StatusOK, response)
}
