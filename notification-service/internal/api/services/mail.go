package services

import (
	"fmt"
	"log"
	"os"
	"time"

	"aerona.thanhtd.com/notification-service/internal/configs"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	cfg *configs.EmailConfig
}

type RetryConfig struct {
	MaxRetries    int
	InitialDelay  time.Duration
	BackoffFactor float64
}

func NewEmailService(cfg *configs.EmailConfig) *EmailService {
	return &EmailService{cfg: cfg}
}

func (s *EmailService) sendWithRetry(msg *gomail.Message, retryCfg RetryConfig) error {
	var err error
	currentDelay := retryCfg.InitialDelay

	for attempt := 1; attempt <= retryCfg.MaxRetries+1; attempt++ {
		err = s.cfg.Dialer.DialAndSend(msg)
		if err == nil {
			log.Printf("Email sent successfully on attempt %d", attempt)
			return nil
		}

		if attempt == retryCfg.MaxRetries+1 {
			return fmt.Errorf("failed to send email after %d attempts: %w", retryCfg.MaxRetries+1, err)
		}

		log.Printf("Attempt %d failed: %v. Retrying in %v...", attempt, err, currentDelay)

		time.Sleep(currentDelay)

		currentDelay = time.Duration(float64(currentDelay) * retryCfg.BackoffFactor)
	}

	return err
}

func (s *EmailService) SendEmail(to, subject, body string, isHTML bool) error {
	if to == "" {
		return fmt.Errorf("recipient email address is required")
	}
	if subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if body == "" {
		return fmt.Errorf("email body is required")
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", s.cfg.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	if isHTML {
		msg.SetBody("text/html", body)
	} else {
		msg.SetBody("text/plain", body)
	}

	retryCfg := RetryConfig{
		MaxRetries:    3,
		InitialDelay:  2 * time.Second,
		BackoffFactor: 2.0,
	}

	err := s.sendWithRetry(msg, retryCfg)
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

func (s *EmailService) SendEmailWithAttachment(to, subject, body, attachmentPath string, isHTML bool) error {
	if to == "" {
		return fmt.Errorf("recipient email address is required")
	}
	if subject == "" {
		return fmt.Errorf("email subject is required")
	}
	if body == "" {
		return fmt.Errorf("email body is required")
	}
	if attachmentPath == "" {
		return fmt.Errorf("attachment path is required")
	}

	if _, err := os.Stat(attachmentPath); os.IsNotExist(err) {
		return fmt.Errorf("attachment file does not exist: %s", attachmentPath)
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", s.cfg.From)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	if isHTML {
		msg.SetBody("text/html", body)
	} else {
		msg.SetBody("text/plain", body)
	}

	msg.Attach(attachmentPath)

	retryCfg := RetryConfig{
		MaxRetries:    3,
		InitialDelay:  2 * time.Second,
		BackoffFactor: 2.0,
	}

	err := s.sendWithRetry(msg, retryCfg)
	if err != nil {
		return fmt.Errorf("failed to send email with attachment to %s: %w", to, err)
	}

	log.Printf("Email with attachment sent successfully to %s", to)
	return nil
}
