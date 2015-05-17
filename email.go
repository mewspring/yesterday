package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/smtp"
	"os"
	"path/filepath"
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

	// Authenticate with the SMTP server.
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
Content-Type: text/plain; charset=utf-8

%s
`
	fmt.Fprintf(buf, format[1:], date, e.to, e.subject, e.message)
	for name, content := range e.attachments {
		l72 := &lineBreaker{
			x: 72,
			w: buf,
		}
		enc := base64.NewEncoder(base64.StdEncoding, l72)
		ext := filepath.Ext(name)
		const format = `--BOUNDARY
Content-Type: %s
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename=%q
`
		fmt.Fprintf(buf, format, mime.TypeByExtension(ext), name)
		enc.Write(content)
		enc.Close()
		buf.WriteByte('\n')
	}
	fmt.Fprintln(buf, "--BOUNDARY--")

	// Optional debug output.
	if flagDebug {
		log.Print("### message in wire format ###\n", buf.String())
	}

	// Send email.
	addr := fmt.Sprintf("%s:%d", auth.Host, auth.Port)
	if err := smtp.SendMail(addr, a, "", to, buf.Bytes()); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// lineBreaker is a io.Writer which inserts a newline after every x bytes
// written.
type lineBreaker struct {
	// Insert a newline after every x bytes written.
	x int
	// Number of bytes written.
	n int
	// Underlying io.Writer.
	w io.Writer
}

// Write writes buf to the underlying io.Writer and injects a newline after
// every x bytes written.
func (l *lineBreaker) Write(buf []byte) (n int, err error) {
	newline := []byte{'\n'}
	for i := 0; i < len(buf); {
		if l.n > 0 && l.n%l.x == 0 {
			// Inject newline.
			if _, err = l.w.Write(newline); err != nil {
				return n, errutil.Err(err)
			}
		}
		end := i + l.x - (l.n % l.x)
		if end > len(buf) {
			end = len(buf)
		}
		m, err := l.w.Write(buf[i:end])
		n += m
		l.n += m
		i += m
		if err != nil {
			return n, errutil.Err(err)
		}
	}
	return n, nil
}
