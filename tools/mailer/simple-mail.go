package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gookit/ini/v2"
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"
)

type userAccount struct {
	Title string
	Url   string
}

type SMTP interface {
	Send(to string, subject string, attachmentPatch string) error
	AuthMail(to, codeConfirmation string, isConfirmation bool) error
}

func AuthMail(to, codeConfirmation string, isConfirmationReset bool) error {
	var templatePath, content string
	ini.LoadFiles("configs/config.ini")
	linkConfirmation := fmt.Sprintf("http://%s:%d/auth/confirm?code=%s", ini.String("domain"), ini.Int("port"), codeConfirmation)
	if codeConfirmation == "" {
		return errors.New("codeConfirmation is not define")
	}
	if isConfirmationReset {
		templatePath = "web/mail/paste-templates/password-reset.html"
		content = "Reset Your Password"
	} else {
		templatePath = "web/mail/paste-templates/email-confirmation.html"
		content = "Confirm Your Email Address"
	}
	data := userAccount{
		Title: content,
		Url:   linkConfirmation,
	}
	var body bytes.Buffer
	var t, err = template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	t.Execute(&body, data)
	err = Send(to, content, body.Bytes(), "", "")
	if err != nil {
		return err
	}

	return nil
}

func Send(to string, subject string, body []byte, attachmentPatch string, filename string) error {
	server := mail.NewSMTPClient()
	server.Host = os.Getenv("SMTP_SERVER")
	server.Port, _ = strconv.Atoi(os.Getenv("SMTP_PORT1"))
	server.Username = os.Getenv("SMTP_LOGIN")
	server.Password = os.Getenv("SMTP_PASSWORD")
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.Authentication = mail.AuthPlain
	server.KeepAlive = false
	if server.Port == 587 {
		server.Encryption = mail.EncryptionSTARTTLS
	} else {
		server.Encryption = mail.EncryptionSSLTLS
	}
	smtpClient, err := server.Connect()
	if err != nil {
		log.Println("server.Connect|", err)
		return err
	}
	email := mail.NewMSG()
	email.SetFrom(os.Getenv("MAIL_NAME")).
		AddTo(to).
		SetSubject(subject).
		SetListUnsubscribe("<mailto:vasek625@gmail.com?subject=https://example.com/unsubscribe>")

	email.SetBodyData(mail.TextHTML, body)
	//email.AddAlternative(mail.TextPlain, message)

	// also you can add body from []byte with SetBodyData, example:
	// email.SetBodyData(mail.TextHTML, []byte(htmlBody))
	// or alternative part
	// email.AddAlternativeData(mail.TextHTML, []byte(htmlBody))

	// add inline
	if attachmentPatch != "" && filename != "" {
		email.Attach(&mail.File{FilePath: attachmentPatch, Name: filename, Inline: false})
	}
	// always check error after send
	if email.Error != nil {
		log.Println("Ошибка в формировании письма ", email.Error)
		return email.Error
	}
	//Call Send and pass the client
	err = email.Send(smtpClient)
	if err != nil {
		return err
	} else {
		return nil
	}
	return nil //
}
