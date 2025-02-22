package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ContactForm represents the expected payload.
type ContactForm struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// sendEmailHandler handles the API Gateway request.
func sendEmailHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Only allow POST requests.
	print("TEST")
	if req.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method not allowed",
		}, nil
	}

	// Check the Origin header to ensure the request is coming from your website.
	allowedOrigin := "https://yourwebsite.com" // Replace with your actual domain.
	origin := req.Headers["origin"]
	if origin != allowedOrigin {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized origin",
		}, nil
	}

	// Parse the request body.
	var form ContactForm
	if err := json.Unmarshal([]byte(req.Body), &form); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid JSON data",
		}, nil
	}

	// Retrieve Mailjet credentials and sender email from environment variables.
	mailjetApiKey := os.Getenv("MAILJET_API_KEY")
	mailjetSecretKey := os.Getenv("MAILJET_SECRET_KEY")
	senderEmail := os.Getenv("SENDER_EMAIL")
	if mailjetApiKey == "" || mailjetSecretKey == "" || senderEmail == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Server configuration error",
		}, nil
	}

	// Build the Mailjet message payload.
	message := struct {
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
	}{}
	message.From.Email = senderEmail
	message.From.Name = "Your Business Name" // Update as needed.
	message.To = []struct {
		Email string `json:"Email"`
		Name  string `json:"Name"`
	}{
		{Email: "destination@example.com", Name: "Recipient Name"}, // Update recipient.
	}
	message.Subject = "New Contact Form Submission"
	message.TextPart = fmt.Sprintf("Email: %s\nPhone: %s\nMessage: %s", form.Email, form.Phone, form.Message)

	// Build the complete request payload.
	requestBody := struct {
		Messages []interface{} `json:"Messages"`
	}{
		Messages: []interface{}{message},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error creating request payload",
		}, nil
	}

	// Create the HTTP request to Mailjet's API.
	mailjetURL := "https://api.mailjet.com/v3.1/send"
	reqMailjet, err := http.NewRequest("POST", mailjetURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error creating Mailjet request",
		}, nil
	}
	reqMailjet.Header.Set("Content-Type", "application/json")
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", mailjetApiKey, mailjetSecretKey)))
	reqMailjet.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(reqMailjet)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error sending request to Mailjet",
		}, nil
	}
	defer resp.Body.Close()

	// Handle errors from Mailjet.
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return events.APIGatewayProxyResponse{
			StatusCode: resp.StatusCode,
			Body:       fmt.Sprintf("Mailjet error: %s", string(bodyBytes)),
		}, nil
	}

	// Return a successful response.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Email sent successfully",
	}, nil
}

func main() {
	lambda.Start(sendEmailHandler)
}
