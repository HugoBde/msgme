package main

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"github.com/withmandala/go-log"
)

// Global logger so we can access it anywhere
var logger *log.Logger

func generateContactFormSubmitHandler() func(http.ResponseWriter, *http.Request) {
	// Static constants
	var app_pwd string
	var username string
	var auth_server string
	var sending_server string
	var sender string
	var recipient string

	var valid bool

	// Read environment variables
	if app_pwd, valid = os.LookupEnv("APPLICATION_PASSWORD"); !valid {
		logger.Fatal("APPLICATION_PASSWORD environment variable not set")
	}

	if username, valid = os.LookupEnv("AUTH_USERNAME"); !valid {
		logger.Fatal("AUTH_USERNAME environment variable not set")
	}

	if auth_server, valid = os.LookupEnv("AUTH_SERVER"); !valid {
		logger.Fatal("AUTH_SERVER environment variable not set")
	}

	if sending_server, valid = os.LookupEnv("SENDING_SERVER"); !valid {
		logger.Fatal("SENDING_SERVER environment variable not set")
	}

	if sender, valid = os.LookupEnv("SENDER"); !valid {
		logger.Fatal("SENDER environment variable not set")
	}

	if recipient, valid = os.LookupEnv("RECIPIENT"); !valid {
		logger.Fatal("RECIPIENT environment variable not set")
	}

	// Build our authentication struct
	auth := smtp.PlainAuth("", username, app_pwd, auth_server)

	// The actual route handler
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request form data into r.Form (GET) and r.PostForm (POST)
		if err := r.ParseForm(); err != nil {
			logger.Error(err)
			return
		}

		// Build message from form data
		msg := []byte(fmt.Sprintf("Subject: Website Enquiry Form\r\nFrom: HugoBde.io\r\nTo: bouderliqueh@gmail.com\r\n\r\nBody: %s\r\n\r\nContact Email: %s\r\nContact Phone Number: %s",
			r.FormValue("msg_body"),
			r.FormValue("contact_email"),
			r.FormValue("contact_phone_no")))

		// Send Email to myself
		err := smtp.SendMail(sending_server, auth, sender, []string{recipient}, msg)

		if err != nil {
			logger.Error(err)
		}
	}
}

func main() {
	// Initialise logger
	logger = log.New(os.Stderr)

	// Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Failed to load .env file")
	}

	// Generate form submission handler
	contactFormSubmitHandler := generateContactFormSubmitHandler()

	// Map handler to /contact_form
	http.HandleFunc("/contact_form", contactFormSubmitHandler)

	// Start listening on port 8080
	logger.Fatal(http.ListenAndServe(":8080", nil))
}
