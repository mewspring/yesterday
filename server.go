package main

import (
	"html/template"
	"io/ioutil"
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
	switch req.Method {
	case "GET":
		if err := srv.serveGET(w, req); err != nil {
			log.Print(err)
		}
	case "POST":
		if err := srv.servePOST(w, req); err != nil {
			log.Print(err)
		}
	default:
		log.Printf("method %q not yet supported", req.Method)
	}
}

// serveGET responds to the given HTTP GET request.
func (srv *emailServer) serveGET(w http.ResponseWriter, req *http.Request) error {
	// Display form page.
	t, err := template.ParseFiles("data/index.html")
	if err != nil {
		return errutil.Err(err)
	}
	if err = t.Execute(w, nil); err != nil {
		return errutil.Err(err)
	}
	return nil
}

// servePOST responds to the given HTTP POST request.
func (srv *emailServer) servePOST(w http.ResponseWriter, req *http.Request) error {
	const MB = 1024 * 1024
	if err := req.ParseMultipartForm(50 * MB); err != nil {
		return errutil.Err(err)
	}
	attachments := make(map[string][]byte, len(req.MultipartForm.File))
	for _, file := range req.MultipartForm.File["attachment"] {
		name := file.Filename
		f, err := file.Open()
		if err != nil {
			return errutil.Err(err)
		}
		buf, err := ioutil.ReadAll(f)
		f.Close()
		if err != nil {
			return errutil.Err(err)
		}
		if _, ok := attachments[name]; ok {
			return errutil.Newf("unable to add attachment; duplicate file name %q", name)
		}
		attachments[name] = buf
	}

	// Send email.
	date := time.Now().Add(-time.Hour * 24)
	e := &Email{
		to:          req.FormValue("to"),
		subject:     req.FormValue("subject"),
		message:     req.FormValue("message"),
		date:        date,
		attachments: attachments,
	}
	if err := e.Send(srv.auth); err != nil {
		return errutil.Err(err)
	}

	// Display success page :)
	t, err := template.ParseFiles("data/enjoy.html")
	if err != nil {
		return errutil.Err(err)
	}
	if err := t.Execute(w, nil); err != nil {
		return errutil.Err(err)
	}
	return nil

}
