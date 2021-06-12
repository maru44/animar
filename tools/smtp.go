package tools

import (
	"animar/v1/configs"
	"fmt"
	"net/smtp"
)

func BaseSendMail(subject string, message string, to string) error {
	smtpAuth := smtp.PlainAuth(
		"",
		configs.EmailHostUser,
		configs.EmailHostPass,
		configs.EmailHost,
	)

	return smtp.SendMail(
		configs.EmailHost+":"+configs.EmailPort,
		smtpAuth,
		configs.EmailHostUser,
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
