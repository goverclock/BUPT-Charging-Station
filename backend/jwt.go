package main

import (
	"buptcs/data"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	UserName string `json:"user_name"`
	jwt.StandardClaims
}

var jwtKey = []byte("my-secret-key")

// 进行 JWT 鉴权
func auth_middleware(c *gin.Context) {
	if !data.JWT_enable {
		return
	}
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{"message": "Authorization header is missing"})
		c.Abort()
		return
	}

	// 解析 JWT
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(401, gin.H{"message": "Invalid token signature"})
			c.Abort()
			return
		}
		c.JSON(401, gin.H{"message": "Invalid token"})
		c.Abort()
		return
	}

	// 验证 JWT 是否过期
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		c.JSON(401, gin.H{"message": "Invalid token"})
		c.Abort()
		return
	}

	// 将用户 ID 保存到上下文中
	c.Set("user_name", claims.UserName)
	c.Next()
}
