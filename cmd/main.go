package main

import (
	cors2 "GetHotWord/config/cors"
	_ "GetHotWord/interval/api/models"
	"GetHotWord/interval/api/routes"
	"GetHotWord/logger"
	"GetHotWord/middleware"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// gin日志记录的地方
	// file, ferr := os.OpenFile("../log/log.log", os.O_RDWR|os.O_CREATE, 0755)
	// if ferr != nil {
	// 	fmt.Println(ferr)
	// }
	// gin.DefaultWriter = io.MultiWriter(file, os.Stdout)
	logger.Init()
	r := gin.New()

	r.Use(cors.New(cors2.GetCors()))
	r.Use(middleware.ErrHandler())
	r.Use(middleware.GinLogger(zap.L()), middleware.GinRecovery(zap.L(), true))
	// 登录拦截器
	r.Use(middleware.TokenHanlder())
	r.MaxMultipartMemory = 8 << 20

	//跨域设置
	routes.Routes(r)
	// err := r.Run(":8080")
	err := r.RunTLS(":8080", "cert.pem", "key.pem")
	if err != nil {
		fmt.Println("启动出错了")
		return
	}

}
