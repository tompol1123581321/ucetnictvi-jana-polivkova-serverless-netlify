package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// ContactForm represents the expected payload.
type ContactForm struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func sendEmailHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Log startup message
	fmt.Println("Spinning up the lambda")

	// Allow only POST requests.
	if req.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method not allowed",
		}, nil
	}

	// Load configuration.
	cfg := loadConfig()

	// Check that the request is coming from an allowed origin.
	origin := req.Headers["origin"]
	if origin != cfg.AllowedOrigin {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized origin",
		}, nil
	}

	// Decode the request body.
	var form ContactForm
	if err := json.Unmarshal([]byte(req.Body), &form); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid JSON data",
		}, nil
	}

	// Send the email via Mailjet.
	result, err := sendMailjetEmail(form, cfg)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error sending email: %v", err),
		}, nil
	}

	// Return a successful response.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       result,
	}, nil
}
