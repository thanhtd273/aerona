package configs

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Dialer   *gomail.Dialer
}

func LoadEmailConfig() (*EmailConfig, error) {

	portStr := os.Getenv("EMAIL_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid EMAIL_PORT: %v", err)
	}

	cfg := EmailConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     port,
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		From:     os.Getenv("EMAIL_FROM"),
	}

	// Validate configuration
	if cfg.Host == "" || cfg.Port == 0 || cfg.Username == "" || cfg.Password == "" || cfg.From == "" {
		return nil, fmt.Errorf("missing required email configuration")
	}

	// Create a new dialer for Gmail SMTP
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	cfg.Dialer = dialer

	return &cfg, nil
}
