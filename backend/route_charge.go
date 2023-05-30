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
		User_id      int     `json:"user_id"` // unused
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
		} `json:"data"`
	}

	// get user
	user_name, _ := ctx.Get("user_name") // ctx.Set in JWT
	user, err := data.UserByName(user_name.(string))
	if err != nil {
		log.Fatal("UserByName")
	}

	// create car for the user
	car := data.Car{
		OwnedBy: user.Uuid,
	}
	car.ChargeMode = request.ChargeMode
	car.ChargeAmount = request.ChargeAmount

	if car.ChargeAmount == 0.0 {
		response.Code = CodeKeyError
		response.Msg = "charge amount should not be 0"
	} else if !scheduler.JoinCar(user, &car) { // try to join the car in the waiting area
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
	if err != nil {
		log.Fatal("UserByName")
	}
	car, err := scheduler.CarByUser(user)
	if err != nil {
		// user hasn't submit the car
		log.Println(err)
		response.Code = CodeForbidden
		response.Msg = "user hasn't submit charge"
	} else {
		response.Code = CodeSucceed
		response.Msg = "succeed"
		wc, err := scheduler.WaitCountByCar(car)
		if err != nil {
			log.Fatal(err)
		}
		response.Data.Waiting_count = wc
	}

	ctx.JSON(http.StatusOK, response)
}

func charge_details(ctx *gin.Context) {
	var request struct {
		User_id int `json:"user_id"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int           `json:"code"`
		Msg  string        `json:"msg"`
		Data []data.Report `json:"data"`
	}

	user, err := data.UserById(request.User_id)
	if err != nil {
		log.Fatal(err)
	}
	rps := scheduler.ReportsByUser(user)
	response.Code = CodeSucceed
	response.Msg = "succeed"
	response.Data = rps

	ctx.JSON(http.StatusOK, response)
}

func charge_startCharge(ctx *gin.Context) {
	var request struct {
		User_id int `json:"user_id"`
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
		log.Fatal("ctx.Get(user_name)")
	}
	user, err := data.UserByName(user_name.(string))
	if err != nil {
		log.Fatal(err, " no such user ", request.User_id)
	}
	car, err := scheduler.CarByUser(user)
	if err != nil {
		response.Code = CodeForbidden
		response.Msg = "user hasn't submit charge"
	} else {
		err = scheduler.StartChargeCar(car)
		if err != nil {
			response.Code = CodeForbidden
			response.Msg = "car is not ready to charge"
		} else { // ok to start charge
			response.Code = CodeSucceed
			response.Msg = "charge started"
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func charge_cancelCharge(ctx *gin.Context) {
	var request struct {
		User_id int `json:"user_id"`
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
		log.Fatal("ctx.Get")
	}
	user, err := data.UserByName(user_name.(string))
	if err != nil {
		log.Fatal("UserByName: ", user_name)
	}

	if scheduler.CancelCharge(user) {
		response.Code = CodeSucceed
		response.Msg = "cancel succeeded"
	} else {
		response.Code = CodeForbidden
		response.Msg = "user is charging, should end charge"
	}

	ctx.JSON(http.StatusOK, response)
}
