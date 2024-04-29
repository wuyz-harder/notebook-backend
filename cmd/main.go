package main

import (
	"fmt"

	cors2 "github.com/wuyz-harder/notebook-backend/config/cors"
	_ "github.com/wuyz-harder/notebook-backend/interval/api/models"
	"github.com/wuyz-harder/notebook-backend/interval/api/routes"
	"github.com/wuyz-harder/notebook-backend/logger"
	"github.com/wuyz-harder/notebook-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	logger.Init()
	r := gin.New()
	//  中间件
	r.Use(cors.New(cors2.GetCors()))
	// 错误拦截器
	r.Use(middleware.ErrHandler())
	r.Use(middleware.GinLogger(zap.L()), middleware.GinRecovery(zap.L(), true))
	// 登录拦截器
	r.Use(middleware.TokenHanlder())
	// 静态路径
	r.Static("/file", "../file")

	r.MaxMultipartMemory = 8 << 20
	//跨域设置
	routes.Routes(r)
	err := r.Run(":8080")
	// 下面是https
	// err := r.RunTLS(":8080", "cert.pem", "key.pem")
	if err != nil {
		fmt.Println("启动出错了")
		return
	}

}
