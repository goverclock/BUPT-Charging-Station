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
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ctx.BindJSON(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			User_id   int    `json:"user_id"`
			User_type int    `json:"user_type"` // 0 - regular user, 1 - admin
			Token     string `json:"token"`     // may be used later
		} `json:"data"`
	}

	// authenticate
	user, err := data.UserByName(request.Username)
	response.Data.User_id = user.Id
	if user.IsAdmin {
		response.Data.User_type = 1
	} else {
		response.Data.User_type = 0
	}
	// if user.Id == 1 { // note: in our database, user with id == 1 is considered admin
	// 	response.Data.User_type = 1
	// } else {
	// 	response.Data.User_type = 0
	// }
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
	if response.Code == CodeSucceed {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			UserName: request.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		})
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			ctx.JSON(500, gin.H{"message": "Internal server error"})
			return
		}

		response.Data.Token = tokenString
		ctx.Header("Authorization", tokenString)

	}
	ctx.JSON(http.StatusOK, response)
}

func register_user(ctx *gin.Context) {
	amazing_lock.Lock()
	defer amazing_lock.Unlock()
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ctx.BindJSON(&request)
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			User_id int `json:"user_id"`
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
		err = user.Create(false) // save user register information
		if err != nil {
			log.Fatal(err, "fail to create user")
		}

		response.Data.User_id = user.Id
	}

	ctx.JSON(http.StatusOK, response)
}
