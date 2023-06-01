package main

import (
	"buptcs/data"
	"buptcs/scheduler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func system_getsettings(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()

	// no body from request
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Waiting_area_size  int `json:"waiting_area_size"`
			Charging_queue_len int `json:"charging_queue_len"`
			Call_schedule      int `json:"call_schedule"`
			Fault_schedule     int `json:"fault_schedule"`
		} `json:"data"`
	}

	response.Code = CodeSucceed
	response.Msg = "succeeded"
	response.Data.Waiting_area_size = data.MAX_WAITING_SLOT
	response.Data.Charging_queue_len = data.MAX_STATION_QUEUE
	response.Data.Call_schedule = data.CALL_SCHEDULE
	response.Data.Fault_schedule = scheduler.GetFautlSchedule()

	ctx.JSON(http.StatusOK, response)
}

func system_setsettings(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()

	var request struct {
		Waiting_area_size  int `json:"waiting_area_size"`
		Charging_queue_len int `json:"charging_queue_len"`
		Call_schedule      int `json:"call_schedule"`
		Fault_schedule     int `json:"fault_schedule"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
		} `json:"data"`
	}

	scheduler.ChangeSettings(request.Waiting_area_size, request.Charging_queue_len, request.Call_schedule, request.Fault_schedule)
	response.Code = CodeSucceed
	response.Msg = "succeeded"

	ctx.JSON(http.StatusOK, response)
}
