package cors

import (
	"time"

	"github.com/gin-contrib/cors"
)

func GetCors() cors.Config {
	cor := cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"PUT", "PATCH", "POST", "OPTIONS", "GET", "DELETE"},
		AllowHeaders: []string{"*"},

		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cor
}
