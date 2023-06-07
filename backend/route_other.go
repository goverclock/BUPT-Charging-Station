package main

import (
	"buptcs/scheduler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getuserinfo(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		UserId int `json:"userid"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			User_id         int     `json:"user_id"`
			Username        string  `json:"username"`
			Tell            string  `json:"tell"`
			Car_battery_now float64 `json:"car_battery_now"`
			Car_battery     float64 `json:"car_battery"`
		} `json:"data"`
	}

	response.Code = CodeSucceed
	response.Msg = "成功获取用户信息"
	response.Data.User_id = request.UserId
	user := scheduler.UserByContext(ctx, request.UserId)
	response.Data.Username = user.Name
	// response.Data.Tell = nil
	response.Data.Car_battery_now = 0
	response.Data.Car_battery = user.BatteryCapacity

	ctx.JSON(http.StatusOK, response)
}

func recharge(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		Recharge_amount float64 `json:"recharge_amount"`
		User_id         int     `json:"user_id"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
		} `json:"data"`
	}

	user := scheduler.UserByContext(ctx, request.User_id)
	if request.Recharge_amount < 0 {
		response.Code = CodeKeyError
		response.Msg = "充值金额必须不小于0"
	} else {
		user.Balance += request.Recharge_amount
		user.Update()
		response.Code = CodeSucceed
		response.Msg = "充值成功"
	}

	ctx.JSON(http.StatusOK, response)
}

func getbalance(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		User_id int `json:"user_id"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Balance float64 `json:"balance"`
		} `json:"data"`
	}

	user := scheduler.UserByContext(ctx, request.User_id)
	response.Code = CodeSucceed
	response.Msg = "查询成功"
	response.Data.Balance = user.Balance

	ctx.JSON(http.StatusOK, response)
}
