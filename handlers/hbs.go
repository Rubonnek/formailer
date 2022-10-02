package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/djatwood/formailer"
	"github.com/djatwood/formailer/logger"
	"github.com/google/martian/log"
)

func hbsResponse(w http.ResponseWriter, code int, err error) {
	r := response{Ok: true}
	if err != nil {
		r.Ok = false
		r.Error = err.Error()
		w.Header().Set("location", w.Header().Get("location")+"?error="+err.Error())
		logger.Error(err)
	}

	body, err := json.Marshal(r)
	if err != nil {
		log.Errorf("failed to marshal response: %w", err)
	}

	w.WriteHeader(code)
	w.Write(body)
}

// Quick hack to redirect users when the contact form has been submitted successfully or otherwise
var succesSlug string = "/contact/success"
var failureSlug string = "/contact/failure"

// hbs just needs a normal http handler
func Hbs(c formailer.Config, w http.ResponseWriter, r *http.Request) {
	// Only process POST methods
	if r.Method != "POST" {
		hbsResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}

	// Read the POST request body
	body := new(strings.Builder)
	_, err := io.Copy(body, r.Body)
	if err != nil {
		hbsResponse(w, http.StatusInternalServerError, err)
		return
	}

	// Debug
	//fmt.Printf("=== Body Print BEGINS ===\n")
	//fmt.Printf("%s\n", body.String())
	//fmt.Printf("=== Body Print ENDS ===")

	// Parse the submission and perform validation checks
	submission, err := c.Parse(r.Header.Get("Content-Type"), body.String())

	if err != nil {
		hbsResponse(w, http.StatusBadRequest, err)
		return
	}

	if submission.Form.ReCAPTCHA {
		v, exists := submission.Values["g-recaptcha-response"].(string)
		if !exists || len(v) < 1 {
			hbsResponse(w, http.StatusBadRequest, nil)
			return
		}

		ok, err := VerifyRecaptcha(v)
		if err != nil {
			err = fmt.Errorf("failed to verify reCAPTCHA: %w", err)
			hbsResponse(w, http.StatusInternalServerError, err)
			return
		}
		if !ok {
			hbsResponse(w, http.StatusBadRequest, nil)
			return
		}

		delete(submission.Values, "g-recaptcha-response")
	}

	// Send the email
	err = submission.Send()
	statusCode := http.StatusOK
	if err != nil {
		// Redirect on error to error page
		err = fmt.Errorf("failed to send email: %w", err)
		logger.Error(err)
		statusCode = http.StatusInternalServerError // setting this in order to show the error in case there is no redirection

		if len(submission.Form.Redirect) > 0 {
			// Write redirect response
			w.Header().Add("Location", submission.Form.Redirect + failureSlug)
			hbsResponse(w, http.StatusSeeOther, nil)
			return
		} else {
			// Write no redirect response:
			hbsResponse(w, statusCode, err)
			return
		}
	}

	// Process HTTP client redirection
	if len(submission.Form.Redirect) > 0 {
		// Write redirect response:
		statusCode = http.StatusSeeOther
		w.Header().Add("Location", submission.Form.Redirect + succesSlug)
		logger.Infof("sent %d emails from %s form and redirected client", len(submission.Form.Emails), submission.Values["_form_name"])
		hbsResponse(w, statusCode, nil)
		return
	} else {
		// Write no redirect response
		hbsResponse(w, statusCode, nil)
		logger.Infof("sent %d emails from %s form and did not redirect client", len(submission.Form.Emails), submission.Values["_form_name"])
	}
}
