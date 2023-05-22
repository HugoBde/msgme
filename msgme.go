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
        msg := []byte(fmt.Sprintf("Subject: Website Enquiry Form\r\nFrom: HugoBde.io\r\nTo: bouderliqueh@gmail.com\r\n\r\nName: %s\r\nContact Email: %s\r\nContact Phone Number: %s\r\n\r\nBody: %s",
			r.FormValue("contact_name"),
			r.FormValue("contact_email"),
			r.FormValue("contact_phone_no"),
			r.FormValue("msg_body")))

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

	var port_num string
	var valid bool
    var cert_file string
    var key_file string 

	if port_num, valid = os.LookupEnv("PORT"); !valid {
		logger.Warn("Missing PORT environment variable. Defaulting to port 3000")
		port_num = "3000"
	}

    if cert_file, valid = os.LookupEnv("CERT_FILE"); !valid {
        logger.Fatal("Missing CERT_FILE environment variable")
    }
    
    if key_file, valid = os.LookupEnv("KEY_FILE"); !valid {
        logger.Fatal("Missing KEY_FILE environment variable")
    }

	// Generate form submission handler
	contactFormSubmitHandler := generateContactFormSubmitHandler()

	// Map handler to /contact_form
	http.HandleFunc("/contact_form", contactFormSubmitHandler)

	// Start listening on port 8080
	logger.Fatal(http.ListenAndServeTLS(":"+port_num, cert_file, key_file, nil))
}
