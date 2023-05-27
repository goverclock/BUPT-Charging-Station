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
		ChargeMode   int     `json:"chargeMode"`
		ChargeAmount float64 `json:"chargeAmount"`
	}
	ctx.Bind(&request)

	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	// get user
	user_name, _ := ctx.Get("user_name") // ctx.Set in JWT
	user, err := data.UserByName(user_name.(string))
	// log.Println(user)
	if err != nil {
		log.Fatal("UserByName")
	}

	car := data.Car{
		OwnedBy: user.Uuid,
	}
	car.ChargeMode = request.ChargeMode
	car.ChargeAmount = request.ChargeAmount

	if !scheduler.JoinCar(car) {
		// no available slot
		response.Code = CodeForbidden
		response.Msg = "the waiting queue is full"
	} else {
		response.Code = CodeSucceed
		response.Msg = "charging request submitted succssfully"
	}

	ctx.JSON(http.StatusOK, response)
}

func charge_getChargingMsg(ctx *gin.Context) {
	var request struct {
		Username string `json:"username"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Waiting_count int `json:"waiting_count"`
		} `json:"data"`
	}

	// get user
	user_name, _ := ctx.Get("user_name") // ctx.Set in JWT
	user, err := data.UserByName(user_name.(string))
	// log.Println(user)
	if err != nil {
		log.Fatal("UserByName")
	}
	car, err := scheduler.CarByUser(&user)	
	if err != nil {
		// user hasn't submit the car
		log.Println(err)
		response.Code = CodeForbidden
		response.Msg = "user hasn't submit charge"
	} else {
		response.Code = CodeSucceed
		response.Msg = "succeed"
		wc, err := scheduler.WaitCountByCar(&car)
		if err != nil {
			log.Fatal(err)
		}
		response.Data.Waiting_count = wc
	}

	ctx.JSON(http.StatusOK, response)
}

// func charge_getChargingMsg(ctx *gin.Context) {
// 	var request struct {
// 		Username string `json:"username"`
// 	}
// 	ctx.Bind(&request)
// 	var response struct {
// 		Code int    `json:"code"`
// 		Msg  string `json:"msg"`
// 		Data struct {
// 			Queue_number  string  `json:"queue_number"`
// 			Waiting_count int     `json:"waiting_count"`
// 			Charge_mode   int     `json:"charge_mode"`
// 			Charge_amount float64 `json:"charge_amount"`
// 			Charge_state  int     `json:"charge_state"`
// 		} `json:"data"`
// 	}

// 	user, err := data.UserByName(request.Username)
// 	if err != nil {
// 		log.Fatal("no such user when getting charging msg")
// 	}
// 	car, err := scheduler.CarByUser(&user)
// 	if err != nil {
// 		// charging request isn't submitted yet
// 		response.Code = CodeSucceed
// 		response.Msg = "succeed"
// 		response.Data.Charge_state = 0
// 	} else {
// 		response.Code = CodeSucceed
// 		response.Msg = "succeed"
// 		response.Data.Queue_number = car.QId
// 		response.Data.Waiting_count = 0 // TODO
// 		response.Data.Charge_mode = car.ChargeMode
// 		response.Data.Charge_amount = car.ChargeAmount
// 		sta := -1
// 		if car.Stage == data.Waiting {
// 			sta = 1
// 		}
// 		if car.Stage == data.Queueing {
// 			sta = 2
// 		} else { // "Charging"
// 			sta = 3
// 		}
// 		response.Data.Charge_state = sta
// 	}

// 	ctx.JSON(http.StatusOK, response)
// }
