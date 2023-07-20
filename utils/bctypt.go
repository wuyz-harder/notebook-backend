package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func GetHashPassword(password string) ([]byte, error) {
	res := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(res, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	} else {
		return hashedPassword, nil
	}
}

func JuegeHashPassworCorrect(hashPassword string, password string) (bool, error) {
	res := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if res == nil {
		return true, nil
	} else {
		return false, res
	}
}
