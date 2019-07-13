package jwt

import (
	"do-mall/pkg/util"
	"github.com/gin-gonic/gin"
	"strings"
	"time"

	"do-mall/pkg/e"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var msg string
		var data interface{}

		code = e.OK
		token := c.GetHeader("Authorization")
		if token == "" {
			code = e.UNAUTHORIZED
		} else {
			if !strings.HasPrefix(token, "Bearer ") {
				code = e.UNAUTHORIZED
			} else {
				token = token[len("Bearer "):]
				claims, err := util.ParseToken(token)
				if err != nil {
					code = e.UNAUTHORIZED
					msg = "Invalid token"
				} else if time.Now().Unix() > claims.ExpiresAt {
					code = e.UNAUTHORIZED
					msg = "Invalid token"
				}
				c.Set("AuthData", claims)
			}

		}

		if code != e.OK {
			if msg == "" {
				msg = e.GetMsg(code)
			}
			c.JSON(code, gin.H{
				"code": code,
				"msg":  msg,
				"data": data,
			})

			c.Abort()
			return
		}
		c.Next()

	}
}
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var msg string
		var data interface{}

		code = e.OK
		token := c.GetHeader("Authorization")
		if token == "" {
			code = e.UNAUTHORIZED
		} else {
			if !strings.HasPrefix(token, "Bearer ") {
				code = e.UNAUTHORIZED
			} else {
				token = token[len("Bearer "):]
				claims, err := util.ParseAdmin(token)
				if err != nil {
					code = e.UNAUTHORIZED
					msg = "Invalid token"
				} else if time.Now().Unix() > claims.ExpiresAt {
					code = e.UNAUTHORIZED
					msg = "Invalid token"
				}
				c.Set("AuthData", claims)
			}

		}

		if code != e.OK {
			if msg == "" {
				msg = e.GetMsg(code)
			}
			c.JSON(code, gin.H{
				"code": code,
				"msg":  msg,
				"data": data,
			})

			c.Abort()
			return
		}
		c.Next()

	}
}
