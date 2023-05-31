package main

import (
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
	} else if st.IsDown {
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
		scheduler.SwitchBrokenStation(request.Charge_id, !st.IsDown)
		if st.IsDown {
			response.Msg = "charging station now is fixed"
		} else {
			response.Msg = "charging station now is broken"
		}
	}

	ctx.JSON(http.StatusOK, response)
}
