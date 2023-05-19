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
		ChargeMode   int  `json:"chargeMode"`
		ChargeAmount float64 `json:"chargeAmount"`
	}
	ctx.Bind(&request)

	var response struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	}

	// get user
	user_name, _ := ctx.Get("user_name")		// ctx.Set in JWT
	user, err := data.UserByName(user_name.(string))
	// log.Println(user)
	if err != nil {
		log.Fatal("UserByName")
	}

	car := data.Car{
		// Id:      1,
		OwnedBy: user.Name,
		// Stage: ,
	}
	car.ChargeMode = request.ChargeMode

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
