package main

import "os"

// Config holds configuration values loaded from environment variables.
type Config struct {
	AllowedOrigin    string
	MailjetAPIKey    string
	MailjetSecretKey string
	SenderEmail      string
	SenderName       string
	RecipientEmail   string
	RecipientName    string
}

// loadConfig reads configuration values from environment variables.
func loadConfig() Config {
	cfg := Config{
		AllowedOrigin:    os.Getenv("ALLOWED_ORIGIN"),
		MailjetAPIKey:    os.Getenv("MAILJET_API_KEY"),
		MailjetSecretKey: os.Getenv("MAILJET_SECRET_KEY"),
		SenderEmail:      os.Getenv("SENDER_EMAIL"),
		SenderName:       os.Getenv("SENDER_NAME"),
		RecipientEmail:   os.Getenv("RECIPIENT_EMAIL"),
		RecipientName:    os.Getenv("RECIPIENT_NAME"),
	}

	// Provide defaults if needed.
	if cfg.AllowedOrigin == "" {
		cfg.AllowedOrigin = "https://yourwebsite.com"
	}
	if cfg.SenderName == "" {
		cfg.SenderName = "Your Business Name"
	}
	if cfg.RecipientName == "" {
		cfg.RecipientName = "Recipient Name"
	}

	return cfg
}
