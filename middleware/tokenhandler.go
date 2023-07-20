package middleware

import (
	"GetHotWord/common"
	"GetHotWord/utils"

	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// GinLogger 接收gin框架默认的日志
func TokenHanlder() gin.HandlerFunc {
	return func(c *gin.Context) {

		path := c.Request.URL.Path
		method := c.Request.Method

		// 判断是否是登录页，登录页就直接跳过
		if !strings.Contains(path, "login") {
			// 如果是options或者wss电话就直接跳过
			if method == "OPTIONS" {
				c.Next()
				return
			}
			// 进行token验证
			token := c.Request.Header.Get("token")
			wsToken := c.Request.Header.Get("Sec-WebSocket-Protocol")
			fmt.Println(token)
			if token == "" && wsToken == "" {
				// c.Error(common.NewError(401, 401, "用户未登录"))
				c.AbortWithError(401, common.NewError(401, 401, "用户未登录"))
				return
			} else {
				// 检查正常的http请求
				if token != "" {
					token, claims, err := utils.ParseToken(token)
					//  如果token不可用就这样吧
					if err != nil || !token.Valid {
						c.Error(common.NewError(401, 200401, "用户未登录或已超时"))
						return
					}

					fmt.Println(claims)
					c.Next()

				} else {
					// 检查ws
					nWsToken, claims, err := utils.ParseToken(wsToken)
					if err != nil || !nWsToken.Valid {
						c.Error(common.NewError(401, 200401, "用户未登录或已超时"))
						return
					}
					fmt.Println(claims)
					// 要返回给客户端，不然报错
					c.Next()

				}

			}

		} else {
			//登录页就放开
			c.Next()
		}

	}
}
