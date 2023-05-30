package main

import (
	"buptcs/data"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getbalance(ctx *gin.Context) {
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
		log.Fatal("no such user ", user_name)
	}
	user, err := data.UserByName(user_name.(string))
	if err != nil {
		log.Fatal(err, " no such user ", request.User_id)
	}
	response.Code = CodeSucceed
	response.Msg = "succeed"
	response.Data.Balance = user.Balance

	ctx.JSON(http.StatusOK, response)
}

func recharge(ctx *gin.Context) {
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

	user_name, ok := ctx.Get("user_name")
	if !ok {
		log.Fatal("ctx.Get()")
	}
	user, err := data.UserByName(user_name.(string))
	if err != nil {
		log.Fatal(err, "no such user: ", request.User_id)
	}
	user.Balance += request.Recharge_amount
	log.Println(user)
	user.Update()
	response.Code = CodeSucceed
	response.Msg = "recharge succeeded"

	ctx.JSON(http.StatusOK, response)
}
