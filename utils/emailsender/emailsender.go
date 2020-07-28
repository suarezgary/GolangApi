package emailsender

import (
	"bytes"
	"html/template"
	"net/smtp"
	"path/filepath"

	"github.com/suarezgary/GolangApi/config"
)

var auth smtp.Auth

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

//NewRequest new request
func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

//SendEmail send email Function
func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)

	addr := config.Cfg().SMTPHost + ":" + config.Cfg().SMTPPort
	auth := smtp.PlainAuth("", config.Cfg().SMTPEmail, config.Cfg().SMTPPassword, config.Cfg().SMTPHost)

	if err := smtp.SendMail(addr, auth, config.Cfg().SMTPEmail, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

// ParseTemplate Template Parser
func (r *Request) ParseTemplate(fileName string, data interface{}) error {
	filePrefix, _ := filepath.Abs("./utils/emailsender/templates/")

	t, err := template.ParseFiles(filePrefix + "/" + fileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

// SendWelcomeEmail send welcome Email
func SendWelcomeEmail(email string, userCompleteName string) {
	emailArray := []string{email}
	welcomeEmailReq := NewRequest(emailArray, "Your Account has been created", "")
	welcomeData := map[string]interface{}{
		"Name": userCompleteName,
	}
	welcomeEmailReq.ParseTemplate("WelcomeEmail.html", welcomeData)
	welcomeEmailReq.SendEmail()
}

// SendForgotEmail send Forgot Email
func SendForgotEmail(email string, userCompleteName string, username string, password string) {
	emailArray := []string{email}
	welcomeEmailReq := NewRequest(emailArray, "Password Restore", "")
	forgotData := map[string]interface{}{
		"Name":     userCompleteName,
		"Username": username,
		"Password": password,
	}
	welcomeEmailReq.ParseTemplate("ForgotPass.html", forgotData)
	welcomeEmailReq.SendEmail()
}

// SendChangePass send Change Pass Email
func SendChangePass(email string, userCompleteName string) {
	emailArray := []string{email}
	welcomeEmailReq := NewRequest(emailArray, "Pasword Changed", "")
	forgotData := map[string]interface{}{
		"Name": userCompleteName,
	}
	welcomeEmailReq.ParseTemplate("ChangePass.html", forgotData)
	welcomeEmailReq.SendEmail()
}
