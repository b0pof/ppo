package mail

import (
	"context"
	"fmt"
	"net/smtp"
	"time"

	"github.com/b0pof/ppo/internal/config"
)

type Usecase struct {
	FromAddress string
	Password    string
	Host        string
	Port        string
}

func New(cfg config.SMTPConfig) *Usecase {
	fmt.Println("cfg:", cfg)

	return &Usecase{
		FromAddress: cfg.FromAddress,
		Password:    cfg.Password,
		Host:        cfg.Host,
		Port:        cfg.Port,
	}
}

func (u *Usecase) SendCode(ctx context.Context, email string, code string) error {
	if u.FromAddress == "" || u.Password == "" {
		return fmt.Errorf("SMTP_FROM or SMTP_PASS not set")
	}

	auth := smtp.PlainAuth("", u.FromAddress, u.Password, u.Host)
	addr := fmt.Sprintf("%s:%s", u.Host, u.Port)

	subject := "Subject: Ваш код подтверждения\r\n"
	body := fmt.Sprintf("Ваш код: %s\r\n", code)

	msg := []byte(
		"To: " + email + "\r\n" +
			"From: " + u.FromAddress + "\r\n" +
			subject +
			"\r\n" +
			body +
			"\r\n" +
			"--\r\n" +
			"Это автоматическое сообщение\r\n",
	)

	errChan := make(chan error, 1)

	go func() {
		errChan <- smtp.SendMail(addr, auth, u.FromAddress, []string{email}, msg)
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(30 * time.Second):
		return fmt.Errorf("SMTP timeout after 30 seconds")
	case <-ctx.Done():
		return ctx.Err()
	}
}
