package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ini/handler"
	"ini/pkg/auth"
	"ini/pkg/errno"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaim, err := auth.ParseRequest(c)
		if err != nil {
			fmt.Println(err)
			handler.SendError(c, errno.ErrTokenInvalid, err.Error())
			//终止函数运行
			c.Abort()
			return
		}

		// 跨越中间件取值
		c.Set("username", userClaim.Username)
		c.Set("expiresAt", userClaim.StandardClaims)

		c.Next()
	}
}
