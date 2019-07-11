package jwt

import (
	"do-mall/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"

	"do-mall/pkg/e"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.GetHeader("Authorization")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			if !strings.HasPrefix(token, "Bearer ") {
				code = e.INVALID_PARAMS
			} else {
				token = token[len("Bearer "):]
				claims, err := util.ParseToken(token)
				if err != nil {
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				} else if time.Now().Unix() > claims.ExpiresAt {
					code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				}
			}

		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
