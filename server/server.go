package server

import (
	"buptcs/entity"

	"github.com/gin-gonic/gin"
)

type Server interface {
	Run(string) error
}

type serverImpl struct {
	server *gin.Engine

	accounts []entity.RegInfo
}

func New() Server {
	si := serverImpl{}
	si.server = gin.Default()

	// set up HTTP route
	si.server.POST("/login", func(ctx *gin.Context) {
		ctx.JSON(200, si.TryLogin(ctx))
	})
	si.server.POST("/register", func(ctx *gin.Context) {
		ctx.JSON(200, si.TryReg(ctx))
	})

	return &si
}

func (s *serverImpl) Run(port string) error {
	return s.server.Run(port)
}

func (s *serverImpl)TryLogin(ctx *gin.Context) entity.LogStatus {
	var ls entity.LogStatus
	var li entity.LogInfo
	ctx.BindJSON(&li)
	
	ls.Success = false
	for _, v := range s.accounts {
		if v.Username == li.Username && v.Passwd == li.Passwd {
			ls.Success = true
		}
	}

	return ls
}

func (s *serverImpl)TryReg(ctx *gin.Context) entity.RegStatus {
	var rs entity.RegStatus
	var ri entity.RegInfo
	ctx.BindJSON(&ri)

	rs.Success = true
	for _, v := range s.accounts {
		if v.Username == ri.Username {
			rs.Success = false
		}
	}
	if rs.Success {
		s.accounts = append(s.accounts, ri)
	}
	
	return rs
}
