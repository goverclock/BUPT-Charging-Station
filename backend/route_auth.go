package main

import (
	"buptcs/data"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func login_user(ctx *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ctx.BindJSON(&request)
	var response struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Username        string  `json:"username"`
		} `json:"data"`
	}

	// authenticate
	response.Data.Username = request.Username
	user, err := data.UserByName(request.Username)
	if request.Username == "" || request.Password == "" {
		response.Code = CodeKeyError
		response.Msg = "need user name or password"
	} else if err != nil {
		response.Code = CodeForbidden
		response.Msg = "no such user"
	} else if user.Password != data.Encrypt(request.Password) {
		response.Code = CodeForbidden
		response.Msg = "wrong password"
	} else if user.Password == data.Encrypt(request.Password) {
		response.Code = CodeSucceed
		response.Msg = "user login succeeded"
	}

	// JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserName: request.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})
    tokenString, err := token.SignedString(JwtKey)
    if err != nil {
        ctx.JSON(500, gin.H{"message": "Internal server error"})
        return
    }

	ctx.Header("Authorization", tokenString)

	ctx.JSON(http.StatusOK, response)
}


func register_user(ctx *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ctx.BindJSON(&request)
	log.Println(request)

	var response struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Username        string  `json:"username"`
			Password        string  `json:"password"`
			Balance         float64 `json:"balance"`
			BatteryCapacity float64 `json:"batteryCapacity"`
		} `json:"data"`
	}

	user, err := data.UserByName(request.Username)
	if request.Username == "" || request.Password == "" {
		response.Code = CodeKeyError
		response.Msg = "need user name or password"
	} else if err == nil {
		// user with same name already exists, can't register
		response.Code = CodeForbidden
		response.Msg = "username is taken"
	} else {
		response.Code = CodeSucceed
		response.Msg = "register succeeds"
		user.Name = request.Username
		user.Password = request.Password
		user.IsAdmin = false
		user.Balance = 0
		user.BatteryCapacity = 0
		err = user.Create(false)	// save user register information
		if err != nil {
			log.Fatal(err, "fail to create user")
		}

		response.Data.Username = user.Name
		response.Data.Password = user.Password
		response.Data.Balance = user.Balance
		response.Data.BatteryCapacity = user.BatteryCapacity
	}

	ctx.JSON(http.StatusOK, response)
}
