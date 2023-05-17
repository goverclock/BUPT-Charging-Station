package main

import (
	"buptcs/data"
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

	// 1. check if there's available slot in waiting area
	// i.e. if there are 6 waiting cars
	cars := data.Cars
	waiting_slot := 6
	for _, c := range cars {
		if c.Stage == "Waiting" {
			waiting_slot--
		}
	}
	if waiting_slot < 0 {
		log.Fatal("impossible waiting slot:", waiting_slot)
	} else if waiting_slot == 0 {
		// no available slot
		response.Code = "500"
		response.Msg = "the waiting queue is full"
	} else {
		response.Code = "200"
		response.Msg = "charging request submitted succssfully"
	}

	// TODO: 2. check if current user already has ongoing charging

	ctx.JSON(http.StatusOK, response)
}
