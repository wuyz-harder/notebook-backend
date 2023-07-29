package utils

import "math/rand"

func GenerateClientCode() string {
	var str = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	code := ""
	for i := 1; i < 10; i++ {
		code = code + string(str[rand.Intn(35)])
	}
	return code
}
