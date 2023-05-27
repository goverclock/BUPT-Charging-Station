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
	
	user, err := data.UserById(request.User_id)
	if err != nil {
		log.Fatal(err, " no such user ", request.User_id)
	}
	response.Code = CodeSucceed
	response.Msg = "succeed"
	response.Data.Balance = user.Balance

	ctx.JSON(http.StatusOK, response)
}
