package main

import (
	"buptcs/data"
	"buptcs/scheduler"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)
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

	user := scheduler.UserByContext(ctx)
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

	user_name, ok := ctx.Get("user_name")
	if !ok {
		log.Println("no such user ", user_name)
	}
	user, err := data.UserByName(user_name.(string))
	if err != nil {
		log.Println(err, " no such user ", request.User_id)
	}
	response.Code = CodeSucceed
	response.Msg = "查询成功"
	response.Data.Balance = user.Balance

	ctx.JSON(http.StatusOK, response)
}
