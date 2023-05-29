package main

import (
	"buptcs/scheduler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func chargeports_turnon(ctx *gin.Context) {
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

	st := scheduler.GetStationById(request.Charge_id)
	if st.Failure {
		response.Code = CodeForbidden
		response.Msg = "the station is in failure"
	} else {
		scheduler.SetStation(request.Charge_id, !st.Running, st.Failure)
		response.Code = CodeSucceed
		if st.Running {
			response.Msg = "station off"
		} else {
			response.Msg = "station on"
		}
	}

	ctx.JSON(http.StatusOK, response)
}
