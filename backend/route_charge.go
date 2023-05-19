package main

import (
	"buptcs/data"
	"buptcs/scheduler"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func charge_submit(ctx *gin.Context) {
	var request struct {
		ChargeMode   string `json:"chargeMode"`
		ChargeAmount float64    `json:"chargeAmount"`
	}
	ctx.Bind(&request)
	log.Println(request)

	var response struct {
		Code string `json:"code"`
		Msg string `json:"msg"`
	}

	car := data.Car{
		Id: 1,
		// OwnedBy: ,
	}
	if !scheduler.JoinCar(car) {
		// no available slot
		response.Code = "500"
		response.Msg = "the waiting queue is full"
	} else {
		response.Code = "200"
		response.Msg = "charging request submitted succssfully"
	}

	ctx.JSON(http.StatusOK, response)
}
