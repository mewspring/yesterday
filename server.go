package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mewkiz/pkg/errutil"
)

// emailServer represents a server capable of sending emails from the past.
type emailServer struct {
	// SMTP authenticator.
	auth *Auth
}

// ServeHTTP responds to the given HTTP request.
func (srv *emailServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		req.ParseForm()
		log.Print("form:", req.Form)
		e := new(Email)
		e.to = req.FormValue("to")
		e.subject = req.FormValue("subject")
		e.message = req.FormValue("message")
		e.date = time.Now().Add(-time.Hour * 24)
		err := e.Send(srv.auth)
		if err != nil {
			log.Print(err)
		}
		t, err := template.ParseFiles("data/enjoy.html")
		if err != nil {
			log.Fatal(errutil.Err(err))
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal(errutil.Err(err))
		}
		return
	}
	t, err := template.ParseFiles("data/index.html")
	if err != nil {
		log.Fatal(errutil.Err(err))
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(errutil.Err(err))
	}
}
