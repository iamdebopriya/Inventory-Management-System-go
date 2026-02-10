package service

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	SMTPHost     string
	SMTPPort     string
	SenderEmail  string
	SenderPasswd string
	Receiver     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SenderEmail:  os.Getenv("SENDER_EMAIL"),
		SenderPasswd: os.Getenv("SENDER_PASSWORD"),
		Receiver:     os.Getenv("RECEIVER_EMAIL"),
	}
}

func (es *EmailService) send(subject, body string) error {
	auth := smtp.PlainAuth("", es.SenderEmail, es.SenderPasswd, es.SMTPHost)

	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s",
		es.Receiver, subject, body)

	addr := fmt.Sprintf("%s:%s", es.SMTPHost, es.SMTPPort)

	return smtp.SendMail(addr, auth, es.SenderEmail, []string{es.Receiver}, []byte(msg))
}

func (es *EmailService) CategoryMail(name string) {
	subject := "Category Created"
	body := "New Category Created: " + name
	go es.send(subject, body)
}

func (es *EmailService) ProductMail(name string) {
	subject := "Product Created"
	body := "New Product Created: " + name
	go es.send(subject, body)
}

func (es *EmailService) OrderMail(product string, qty int) {
	subject := "Order Created"
	body := fmt.Sprintf("Order Created\nProduct: %s\nQuantity: %d", product, qty)
	go es.send(subject, body)
}
