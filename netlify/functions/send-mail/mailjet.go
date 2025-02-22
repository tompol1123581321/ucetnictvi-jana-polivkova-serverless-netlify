package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// MailjetMessage defines the structure for a single message.
type MailjetMessage struct {
	From struct {
		Email string `json:"Email"`
		Name  string `json:"Name"`
	} `json:"From"`
	To []struct {
		Email string `json:"Email"`
		Name  string `json:"Name"`
	} `json:"To"`
	Subject  string `json:"Subject"`
	TextPart string `json:"TextPart"`
}

// MailjetRequestBody wraps the message.
type MailjetRequestBody struct {
	Messages []MailjetMessage `json:"Messages"`
}

func sendMailjetEmail(form ContactForm, cfg Config) (string, error) {
	// Build the Mailjet message.
	var message MailjetMessage
	message.From.Email = cfg.SenderEmail
	message.From.Name = cfg.SenderName
	message.To = []struct {
		Email string `json:"Email"`
		Name  string `json:"Name"`
	}{
		{Email: cfg.RecipientEmail, Name: cfg.RecipientName},
	}
	message.Subject = "New Contact Form Submission"
	message.TextPart = fmt.Sprintf("Email: %s\nPhone: %s\nMessage: %s", form.Email, form.Phone, form.Message)

	requestBody := MailjetRequestBody{
		Messages: []MailjetMessage{message},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error creating request payload: %v", err)
	}

	// Prepare the Mailjet API request.
	mailjetURL := "https://api.mailjet.com/v3.1/send"
	reqMailjet, err := http.NewRequest("POST", mailjetURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating Mailjet request: %v", err)
	}
	reqMailjet.Header.Set("Content-Type", "application/json")
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.MailjetAPIKey, cfg.MailjetSecretKey)))
	reqMailjet.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(reqMailjet)
	if err != nil {
		return "", fmt.Errorf("error sending request to Mailjet: %v", err)
	}
	defer resp.Body.Close()

	// Check for errors returned by Mailjet.
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("mailjet error: %s", string(bodyBytes))
	}

	return "Email sent successfully", nil
}
