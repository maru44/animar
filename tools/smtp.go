package tools

import (
	"fmt"
	"net/smtp"
	"os"
)

// cannot use
// var hostUser = os.Getenv("EMAIL_HOST_USER")
// var hostPass = os.Getenv("EMAIL_HOST_PASS")
// var host = os.Getenv("EMAIL_HOST")
// var port = os.Getenv("EMAIL_PORT")

func BaseSendMail(subject string, message string, to string) error {
	smtpAuth := smtp.PlainAuth(
		"",
		os.Getenv("EMAIL_HOST_USER"),
		os.Getenv("EMAIL_HOST_PASS"),
		os.Getenv("EMAIL_HOST"),
	)

	return smtp.SendMail(
		os.Getenv("EMAIL_HOST")+":"+os.Getenv("EMAIL_PORT"),
		smtpAuth,
		os.Getenv("EMAIL_HOST_USER"),
		[]string{to},
		[]byte(
			"To: "+to+"\r\n"+
				"Subject: "+subject+"\r\n"+
				"Content-Type: text/plain; charset=\"utf-8\""+"\r\n"+ // anti 文字化け
				"\r\n"+
				message,
		),
	)
}

func SendVerifyEmail(to string, link string) error {
	const subject = "loveAni.me 登録メール"
	var message string
	message = fmt.Sprintf(
		`loveAni.meへのご登録ありがとうございます
		
仮登録が完了しました。
以下のリンクをクリックすることで本登録が完了します。

%s
*** Google社の提供するfirebaseへのリンクとなっております。

loveAni.me
		`, link,
	)

	sended := BaseSendMail(subject, message, to)
	return sended
}
