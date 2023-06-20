package mailer

import (
	"fmt"
	"gofit-api/configs"

	gomail "gopkg.in/gomail.v2"
)

func SendOTP(recipient string, otp int) (string, error) {
	// recipient = "libr.libr1711@gmail.com"
	message := fmt.Sprintf("Here's your OTP : %d. Don't share it with anybody.", otp)
	m := gomail.NewMessage()
	m.SetHeader("From", "gofitapi@gmail.com")
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "GoFit OTP Forgot Password")
	m.SetBody("text/plain", message)

	d := gomail.NewDialer(configs.AppConfig.SMTPHost, 587, configs.AppConfig.SMTPUsername, configs.AppConfig.SMTPPassword)

	if err := d.DialAndSend(m); err != nil {
		return "", err
	}
	return fmt.Sprintf("email sended to %s", recipient), nil
}
