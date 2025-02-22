package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

// MailRequest represents the expected JSON payload.
type MailRequest struct {
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Allow only POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the incoming JSON payload.
	var reqData MailRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Initialize the Mailjet client using environment variables.
	mailjetClient := mailjet.NewMailjetClient(
		os.Getenv("MJ_APIKEY_PUBLIC"),
		os.Getenv("MJ_APIKEY_PRIVATE"),
	)

	// Declare recipients as type mailjet.RecipientsV31.
	recipients := mailjet.RecipientsV31{
		{
			Email: os.Getenv("TO_EMAIL"), // Recipient email from environment.
			Name:  "Recipient Name",      // Customize as needed.
		},
	}

	// Construct the email message payload.
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: os.Getenv("FROM_EMAIL"), // Sender email from environment.
				Name:  "Your Name",             // Customize as needed.
			},
			To:       &recipients, // Use pointer to the properly typed recipients variable.
			Subject:  "New message from website",
			TextPart: "Email: " + reqData.Email + "\nPhone: " + reqData.Phone + "\nMessage: " + reqData.Message,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}

	// Send the email using Mailjet.
	response, err := mailjetClient.SendMailV31(&messages)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	log.Printf("Email sent: %v", response)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully"))
}

func main() {
	http.HandleFunc("/", handler)

	// Use the PORT environment variable if set (Netlify may set this).
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
