package main

import (
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

// Send sends the email from the spoofed date.
func (email *Email) Send() error {
	if len(email.from) < 1 {
		return errutil.New("empty sender email address")
	}
	if len(email.to) < 1 {
		return errutil.New("empty recipient email address")
	}
	panic("not yet implemented.")
}
