package main

import (
	"io"
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
	if _, err := io.WriteString(w, index[1:]); err != nil {
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
	if _, err := io.WriteString(w, success[1:]); err != nil {
		return errutil.Err(err)
	}
	return nil
}

const (
	index = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset='utf-8'>
		<title>Yesterday, all my troubles seemed so far away</title>
	</head>
	<body>
		<form action='/' method='POST' enctype='multipart/form-data'>
			<fieldset style='margin: 10% auto; width: 500px'>
				<legend>Yesterday, all my troubles seemed so far away</legend>
				<div style='float: left; width: 130px'>To:</div>
				<div><input type='text' name='to' placeholder='recipient@example.org' style='width: 360px'></div>
				<div style='float: left; width: 130px'>Subject:</div>
				<div><input type='text' name='subject' placeholder='Subject' style='width: 360px'></div>
				<div style='float: left; width: 130px'>Message:</div>
				<div><textarea name='message' placeholder='Message' style='width: 360px; height: 260px;'></textarea></div>
				<div style='float: left; width: 130px'>Attachments:</div>
				<div><input type='file' name='attachment' multiple style='width: 360px'></div>
				<div style='margin-left: 130px'><input type='submit' value='Send email'></div>
			</fieldset>
		</form>
	</body>
</html>
`
	success = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset='utf-8'>
		<title>Oh, I believe in yesterday</title>
	</head>
	<body>
		<fieldset style='margin: 10% auto; width: 500px'>
			<legend>Oh, I believe in yesterday</legend>
			<em>Email successfully sent 24 hours ago.</em>
		</fieldset>
	</body>
</html>
`
)
