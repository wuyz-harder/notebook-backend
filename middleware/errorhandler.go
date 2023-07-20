package middleware

import (
	"GetHotWord/common"

	"github.com/gin-gonic/gin"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.next()是为了让其他handler先走
		c.Next()
		if length := len(c.Errors); length > 0 {
			e := c.Errors[length-1]
			err := e.Err
			if err != nil {
				var Err *(common.ApiError)
				if e, ok := err.(*(common.ApiError)); ok {
					Err = e
				} else if e, ok := err.(error); ok {
					Err = common.OtherError(e.Error())
				} else {
					Err = common.ServerError
				}
				// 记录一个错误的日志
				c.JSON(Err.StatusCode, Err)
				return
			}
		}

	}
}
