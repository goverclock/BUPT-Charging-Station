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
		User_id      int     `json:"User_id"`	// unused
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

// TODO
func charge_details(ctx *gin.Context) {
	var request struct {
		User_id int `json:"user_id"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Num                   int     `json:"num"`
			Charge_id             int     `json:"charge_id"`
			Charge_mode           int     `json:"charge_mode"`
			Username              string  `json:"username"`
			User_id               int     `json:"user_id"`
			Request_charge_amount float64 `json:"request_charge_amount"`
			Real_charge_amount    float64 `json:"real_charge_amount"`
			Charge_time           int     `json:"charge_time"`
			Charge_fee            float64 `json:"charge_fee"`
			Service_fee           float64 `json:"service_fee"`
			Tot_fee               float64 `json:"tot_fee"`
			Step                  int     `json:"step"`
			Queue_number          string  `json:"queue_number"`
			Subtime               int     `json:"subtime"`
			Inlinetime            int     `json:"inlinetime"`
			Calltime              int     `json:"calltime"`
			Charge_start_time     int     `json:"charge_start_time"`
			Charge_end_time       int     `json:"charge_end_time"`
			Terminate_flag        bool    `json:"terminate_flag"`
			Terminate_time        int     `json:"terminate_time"`
			Failed_flag           bool    `json:"failed_flag"`
			Failed_msg            string  `json:"failed_msg"`
		} `json:"data"`
	}

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

	user, err := data.UserById(request.User_id)
	if err != nil {
		log.Fatal(err, " no such user ", request.User_id)
	}
	log.Println(user)
	car, err := scheduler.CarByUser(user)
	if err != nil {
		response.Code = CodeForbidden
		response.Msg = "user hasn't submit charge"
	} else {
		err = scheduler.StartChargeCar(car)
		if err != nil {
			response.Code = CodeForbidden
			response.Msg = "car is not ready to charge"
		} else {
			response.Code = CodeSucceed
			response.Msg = "charge started"
		}
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
