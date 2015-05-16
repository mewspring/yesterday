package main

import (
	"encoding/json"
	"net/smtp"
	"os"
	"time"

	"github.com/mewkiz/pkg/errutil"
)

// Email represents an email of the past.
type Email struct {
	// Sender email address.
	from string
	// Recipient email address.
	to string
	// Email subject.
	subject string
	// Email message.
	message string
	// Spoof date.
	date time.Time
	// Attachments.
	attachments map[string][]byte
}

// parseAuth parses the provided JSON file and returns an SMTP authenticator.
func parseAuth(path string) (smtp.Auth, error) {
	// Auth specifies the information required for a PLAIN SMTP authenticator.
	type Auth struct {
		// User name.
		User string
		// Password.
		Pass string
		// SMTP host address.
		Host string
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, errutil.Err(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	auth := new(Auth)
	if err := dec.Decode(auth); err != nil {
		return nil, errutil.Err(err)
	}
	return smtp.PlainAuth("", auth.User, auth.Pass, auth.Host), nil
}

// Send sends the email from the spoofed date.
func (email *Email) Send(auth smtp.Auth) error {
	if len(email.from) < 1 {
		return errutil.New("empty sender email address")
	}
	if len(email.to) < 1 {
		return errutil.New("empty recipient email address")
	}
	panic("not yet implemented.")
}
