package main

import "time"

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
	panic("not yet implemented.")
}
