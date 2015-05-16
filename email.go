package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// Auth specifies the information required for a PLAIN SMTP authenticator.
type Auth struct {
	// User name.
	User string
	// Password.
	Pass string
	// SMTP host address.
	Host string
	// Port of SMTP host.
	Port int
}

// parseAuth parses the provided JSON file and returns an SMTP authenticator.
func parseAuth(path string) (*Auth, error) {
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
	return auth, nil
}

// Send sends the email from the spoofed date.
func (e *Email) Send(auth *Auth) error {
	if len(e.from) < 1 {
		return errutil.New("empty sender email address")
	}
	if len(e.to) < 1 {
		return errutil.New("empty recipient email address")
	}
	a := smtp.PlainAuth("", auth.User, auth.Pass, auth.Host)
	to := []string{e.to}
	buf := new(bytes.Buffer)
	const format = `
Date: %s
From: %s
To: %s
Subject: %s
MIME-Version: 1.0
Content-Transfer-Encoding: 8bit
Content-Type: text/html; charset="UTF-8"

%s
`
	date := e.date.Format("Mon, 2 Jan 2006 15:04:05 -0700 (MST)")
	fmt.Fprintf(buf, format[1:], date, e.from, e.to, e.subject, e.message)
	addr := fmt.Sprintf("%s:%d", auth.Host, auth.Port)
	if err := smtp.SendMail(addr, a, e.from, to, buf.Bytes()); err != nil {
		return errutil.Err(err)
	}
	return nil
}
