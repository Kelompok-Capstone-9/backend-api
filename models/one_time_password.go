package models

import (
	"math/rand"
	"time"
)

type OTP struct {
	Email string `json:"email"`
	Code  int    `json:"otp"`
	NewPassword string `json:"new_password"`
}

func (o *OTP) GenerateOTP() error {
	min := 1000
	max := 9000
	// set seed
	rand.NewSource(time.Now().UnixNano())
	// generate random number and print on console
	o.Code = rand.Intn(max-min) + min
	return nil
}
