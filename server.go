package main

import (
	"net/http"
	"net/smtp"
)

// emailServer represents a server capable of sending emails from the past.
type emailServer struct {
	// SMTP authenticator.
	auth smtp.Auth
}

// ServeHTTP responds to the given HTTP request.
func (srv *emailServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	panic("not yet implemented.")
}
