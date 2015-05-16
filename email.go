package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/mewkiz/pkg/errutil"
)

// Email represents an email of the past.
type Email struct {
	// Recipient email address.
	to string
	// Email subject.
	subject string
	// Email message.
	message string
	// Spoof date.
	date time.Time
	// Attachments, as a map from file names to file content.
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
	if len(e.to) < 1 {
		return errutil.New("empty recipient email address")
	}
	a := smtp.PlainAuth("", auth.User, auth.Pass, auth.Host)
	to := []string{e.to}
	buf := new(bytes.Buffer)
	date := e.date.Format("Mon, 2 Jan 2006 15:04:05 -0700 (MST)")

	// Create wire-format message.
	const format = `
Content-Type: multipart/mixed; boundary=BOUNDARY
Date: %s
To: %s
Subject: %s

--BOUNDARY
Content-Type: text/plain; charset="UTF-8"

%s
`
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	fmt.Fprintf(buf, format[1:], date, e.to, e.subject, e.message)
	for name, content := range e.attachments {
		const format = `--BOUNDARY
Content-Type: text/plain
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename=%q; modification-date=%q
`
		fmt.Fprintf(buf, format, name, date)
		// TODO: Split base64 encoded content into line-72.
		enc.Write(content)
		buf.WriteByte('\n')
	}
	fmt.Fprintln(buf, "--BOUNDARY--")

	// TODO: Remove debug output.
	log.Print("### wire-format message ###", buf.String())

	addr := fmt.Sprintf("%s:%d", auth.Host, auth.Port)
	if err := smtp.SendMail(addr, a, "", to, buf.Bytes()); err != nil {
		return errutil.Err(err)
	}
	return nil
}
