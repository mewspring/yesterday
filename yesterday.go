// Yesterday is a procrastination tool which allows you to send emails up to 24
// hours in the past.
//
// It has two modes.
//
// Without the -http flag, it runs in command-line mode and sends email from the
// terminal.
//
//    yesterday -from="john.doe@student.uni.edu" -to="jane.roe@uni.edu" -subject="Report" -message="See attached report." report.pdf
//
// With the -http flag, it runs as a web server and sends email from a web page.
//
//    yesterday -http=:6565
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/mewkiz/pkg/errutil"
)

//go:generate usagen yesterday

// use specifies the usage message of yesterday.
const use = `
Usage: yesterday [OPTION]... FILE...
Send emails up to 24 hours in the past.

Flags:`

// usage prints to standard error a usage message documenting all defined
// command-line flags.
func usage() {
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	// Server mode flags.
	var (
		// flagAddr specifies the HTTP service address (e.g. ":6565").
		flagAddr string
	)

	// Command-line mode flags.
	var (
		// flagFrom specifies the sender email address.
		flagFrom string
		// flagTo specifies the recipient email address.
		flagTo string
		// flagSubject specifies the email subject.
		flagSubject string
		// flagMessage specifies the email message.
		flagMessage string
		// flagPast specifies the spoof date in number of hours in the past.
		flagPast time.Duration
	)

	// Server mode flags.
	flag.StringVar(&flagAddr, "http", "", `HTTP service address (e.g. ":6565").`)

	// Command-line mode flags.
	flag.StringVar(&flagFrom, "from", "", "Sender email address.")
	flag.StringVar(&flagTo, "to", "", "Recipient email address.")
	flag.StringVar(&flagSubject, "subject", "", "Email subject.")
	flag.StringVar(&flagMessage, "message", "", "Email message.")
	flag.DurationVar(&flagPast, "past", 24*time.Hour, "Spoof date in number of hours in the past.")

	flag.Usage = usage

	flag.Parse()

	// Server mode.
	if len(flagAddr) > 0 {
		log.Fatal(listen(flagAddr))
	}

	// Client mode.
	if flagPast < 0 || flagPast > 24*time.Hour {
		log.Fatalf("invalid number of hours in the past; expected >= 0h and <= 24h, got %v", flagPast)
	}
	if len(flagFrom) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	if len(flagTo) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	date := time.Now().Add(-flagPast)
	email := &Email{
		from:        flagFrom,
		to:          flagTo,
		subject:     flagSubject,
		message:     flagMessage,
		date:        date,
		attachments: make(map[string][]byte),
	}
	attachments, err := readAttachments(flag.Args())
	if err != nil {
		log.Fatal(err)
	}
	email.attachments = attachments
	if err := email.Send(); err != nil {
		log.Fatal(err)
	}
}

// readAttachments reads the provided files and returns a mapping from file
// names to file content.
func readAttachments(paths []string) (map[string][]byte, error) {
	attachments := make(map[string][]byte, len(paths))
	for _, path := range paths {
		name := filepath.Base(path)
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errutil.Err(err)
		}
		if _, ok := attachments[name]; ok {
			return nil, errutil.Newf("unable to add attachment %q; duplicate file name %q", path, name)
		}
		attachments[name] = buf
	}
	return attachments, nil
}
