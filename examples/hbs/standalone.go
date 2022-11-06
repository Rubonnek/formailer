package main

// NOTES:
// The main function in Golang is a function with the same name as the package name above.
// In this case, the main function is called "main".
// Also, every package in Golang can have many init functions -- not just one.

import (
	// Formailer
	"github.com/Rubonnek/formailer"
	"github.com/Rubonnek/formailer/handlers"


	// For setting up a local instance we need to setup an http server
	"net/http"

	// For debugging
	//"fmt"
)

// Formailer handles all form submissions
func Formailer(w http.ResponseWriter, r *http.Request) {
	// NOTE: formailer stores a map of strings to Form objects underneath
	// accessible through formailer.DefaultConfig -- this mapping expects
	// the _form_name to be a matching key to a form configuration we have
	// generated here.

	// Let's create a new form configuration and add it to the formailer.DefaultConfig:
	formailer_form := formailer.New("Contact")

	// Debug:
	//fmt.Printf("Default Config: %#v\n", formailer.DefaultConfig) // should be populated with the "contact" form

	// Add a new Email to send this form to:
	formailer_form.AddEmail(formailer.Email{
		// Email ID:
		ID:      "contact",

		// NOTE: The email ID is used to check the SMTP environment variables to see if there is a matching email ID present.
		// Basically email ID allows you to tweak the SMTP relay
		// settings. For example, if we define SMTP_formailer_form_HOST
		// environment variable, then the SMTP host target will be
		// changed for this specific email ID. Otherwise the default one will be used.
		// NOTE: This is really useful for sending emails through different SMTP endpoints or configurations basically.

		// Internally, in email.go specifically, you'll see this (where e.ID is the email ID):
			// prefix := fmt.Sprintf("SMTP_%s_", strings.ToUpper(e.ID))
			// host := or(os.Getenv(prefix+"HOST"), os.Getenv("SMTP_HOST"))
			// user := or(os.Getenv(prefix+"USER"), os.Getenv("SMTP_USER"))
			// pass := or(os.Getenv(prefix+"PASS"), os.Getenv("SMTP_PASS"))
			// defaultPort := os.Getenv("SMTP_PORT")
			// emailPort := os.Getenv(prefix + "PORT")

		// Target email:
		To:      "info@domain.com",

		// Source email:
		From:    `"Company" <noreply@domain.com>`,

		// Email subject:
		Subject: "New Contact Submission",
	})

	// Enable ReCAPTCHA authentication
	//formailer_form.ReCAPTCHA = true

	// When the form handler is successful, redirect to the following page:
	//formailer_form.Redirect = "http://success.com"

	// Debug
	//fmt.Printf("Form: %#v\n", formailer_form)

	// Pass the configuration to the form handler the form handler for the hugo-bootstrap-theme by Razon Yang
	handlers.Hbs(formailer.DefaultConfig, w, r)
}

func debug_http_server() {
	// Dummy server to test Implementation
	http.HandleFunc("/", Formailer)
	http.ListenAndServe(":8080", nil)
}


func main() {
	debug_http_server()
}
