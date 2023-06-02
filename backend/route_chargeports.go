package main

import (
	"buptcs/data"
	"buptcs/scheduler"
	"net/http"

	"github.com/gin-gonic/gin"
)

// on/off
func chargeports_switch(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		Charge_id int `json:"charge_id"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
		} `json:"data"`
	}

	// st is a copy, so no race condition
	st, err := scheduler.StationById(request.Charge_id)
	if err != nil {
		response.Code = CodeKeyError
		response.Msg = "no such charging station"
	} else if st.GetIsDown() {
		response.Code = CodeForbidden
		response.Msg = "the station is in failure"
	} else {
		// turn on/off the station
		scheduler.SwitchStation(request.Charge_id, !st.Running)
		response.Code = CodeSucceed
		if st.Running {
			response.Msg = "station off"
		} else {
			response.Msg = "station on"
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func chargeports_switchBroken(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		Charge_id int `json:"charge_id"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
		} `json:"data"`
	}

	st, err := scheduler.StationById(request.Charge_id)
	if err != nil {
		response.Code = CodeKeyError
		response.Msg = "no such charging station"
	} else {
		response.Code = CodeSucceed
		scheduler.SwitchBrokenStation(request.Charge_id, !st.GetIsDown())
		if st.GetIsDown() {
			response.Msg = "charging station now is broken"
		} else {
			response.Msg = "charging station now is fixed"
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func chargeports_waitingCars(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		Charge_id int `json:"charge_id"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data []struct {
			Username         string  `json:"username"`
			User_id          int     `json:"user_id"`
			Waiting_time     int64   `json:"waiting_time"`
			Charge_amount    float64 `json:"charge_amount"`
			Battery_capacity float64 `json:"battery_capacity"`
		} `json:"data"`
	}

	st, err := scheduler.StationById(request.Charge_id)
	if err != nil {
		response.Code = CodeForbidden
		response.Msg = "no such charging station"
	} else {
		response.Code = CodeSucceed
		response.Msg = "查询成功"
		ent := struct {
			Username         string  `json:"username"`
			User_id          int     `json:"user_id"`
			Waiting_time     int64   `json:"waiting_time"`
			Charge_amount    float64 `json:"charge_amount"`
			Battery_capacity float64 `json:"battery_capacity"`
		}{}

		for _, c := range st.GetQueue() {
			user := data.UserByUUId(c.OwnedBy)
			rp := scheduler.OngoingCopyByUser(user)
			ent.Username = user.Name
			ent.User_id = user.Id
			ent.Waiting_time = (rp.Inlinetime - rp.Subtime) / 60
			ent.Charge_amount = rp.Real_charge_amount
			ent.Battery_capacity = user.BatteryCapacity

			response.Data = append(response.Data, ent)
		}
	}

	ctx.JSON(http.StatusOK, response)
}

// specified range
func chargeports_getreport(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		StartDate int64 `json:"startDate"`
		EndDate   int64 `json:"endDate"`
	}
	ctx.Bind(&request)
	var response struct {
		Code int                  `json:"code"`
		Msg  string               `json:"msg"`
		Data []data.StationReport `json:"data"`
	}

	response.Code = CodeSucceed
	response.Msg = "succeeded"
	response.Data = scheduler.AllStationReports(request.StartDate, request.EndDate)

	ctx.JSON(http.StatusOK, response)
}

// till now
func chargeports_getreports(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var response struct {
		Code int                  `json:"code"`
		Msg  string               `json:"msg"`
		Data []data.StationReport `json:"data"`
	}

	response.Code = CodeSucceed
	response.Msg = "succeeded"
	response.Data = scheduler.AllStationReports(0, 1893456000)

	ctx.JSON(http.StatusOK, response)
}
