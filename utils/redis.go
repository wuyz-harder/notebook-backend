package utils

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		Password: "Aa123456",       // 密码，如果Redis服务器需要密码认证
		DB:       0,                // 选择Redis数据库，默认为0
	})
	// 使用context控制连接操作的超时
	ctx := context.Background()

	// 测试连接是否成功
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
		return
	}
	fmt.Println("Connected to Redis:", pong)

}
