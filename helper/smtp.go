package helper

import (
	"net/smtp"
	"os"
)

var hostUser = os.Getenv("EMAIL_HOST_USER")
var hostPass = os.Getenv("EMAIL_HOST_PASS")
var host = os.Getenv("EMAIL_HOST")
var port = os.Getenv("EMAIL_PORT")

func BaseSendMail(subject string, message string, to string) error {
	smtpAuth := smtp.PlainAuth(
		"",
		hostUser,
		hostPass,
		host,
	)

	return smtp.SendMail(
		host+port,
		smtpAuth,
		hostUser,
		[]string{to},
		[]byte(
			"To: "+to+"\r\n"+
				"Subject:"+subject+"\r\n"+
				"\r\n"+
				message,
		),
	)
}
