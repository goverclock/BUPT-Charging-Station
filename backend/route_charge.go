package main

import (
	"buptcs/data"
	"buptcs/scheduler"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func charge_submit(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		Charge_mode   int     `json:"charge_mode"`
		Charge_amount float64 `json:"charge_amount"`
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
		log.Println("UserByName")
	}

	cp, sp := scheduler.GetFee()
	fee := request.Charge_amount * (cp + sp)
	// check if the user already has a submit
	if scheduler.OngoingCopyByUser(user).Num != 0 {
		response.Code = CodeForbidden
		response.Msg = "已经有进行中的请求了"
	} else if fee > user.Balance{	// if balance not enough, refuse submit
		response.Code = CodeForbidden
		response.Msg = "余额不足"
	} else {
		// create car for the user
		car := data.Car{
			OwnedBy: user.Uuid,
		}
		car.ChargeMode = request.Charge_mode
		car.ChargeAmount = request.Charge_amount

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
}

func charge_getChargingMsg(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
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
		log.Println("UserByName")
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
			log.Println(err)
		}
		response.Data.Waiting_count = wc
	}

	ctx.JSON(http.StatusOK, response)
}

// "chargeSubmit" also goes here
func charge_changeSubmit(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		Charge_mode   int     `json:"charge_mode"`
		Charge_amount float64 `json:"charge_amount"`
		User_id       int     `json:"user_id"` // unused
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
		} `json:"data"`
	}

	// get user
	user := scheduler.UserByContext(ctx)
	if scheduler.ChangeCharge(user, request.Charge_mode, request.Charge_amount) {
		response.Code = CodeSucceed
		response.Msg = "change succeeded"
	} else {
		response.Code = CodeForbidden
		response.Msg = "change failed"
	}

	ctx.JSON(http.StatusOK, response)
}

func charge_cancelCharge(ctx *gin.Context) {
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
		} `json:"data"`
	}

	user := scheduler.UserByContext(ctx)
	if scheduler.CancelCharge(user) {
		response.Code = CodeSucceed
		response.Msg = "cancel succeeded"
	} else {
		response.Code = CodeForbidden
		response.Msg = "user hasn't submitted or is charging, should end charge"
	}

	ctx.JSON(http.StatusOK, response)
}

func charge_startCharge(ctx *gin.Context) {
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
		} `json:"data"`
	}

	user := scheduler.UserByContext(ctx)
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

func charge_end_charge(ctx *gin.Context) {
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
		} `json:"data"`
	}

	user := scheduler.UserByContext(ctx)
	if scheduler.EndCharge(user) {
		response.Code = CodeSucceed
		response.Msg = "end succeeded"
	} else {
		response.Code = CodeForbidden
		response.Msg = "user is not charging"
	}
	ctx.JSON(http.StatusOK, response)
}

func charge_details(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		User_id int `json:"user_id"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int           `json:"code"`
		Msg  string        `json:"msg"`
		Data []data.Report `json:"data"`
	}

	user := scheduler.UserByContext(ctx)
	rps := scheduler.ReportsByUser(user)
	response.Code = CodeSucceed
	response.Msg = "succeed"
	response.Data = rps

	ctx.JSON(http.StatusOK, response)
}
