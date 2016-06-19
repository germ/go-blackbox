package blackbox 

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

/* General usage:
	Using your ID make a call to create or info, this will return a 
	structure with the session ID/all sessions. Save these for future reference

	The structure has a Upload method that will drain reader r and
	store it on my server. Max upload size is 10MB, split larger files
	or host elsewhere and provide a link. Please don't upload more then
	500MB of data, to my server thank you.

	Lastly when all assets have been submitted, call the Finalize method
	to submit it for approval. If I decline the work, an email will be sent
	to the registered account. Otherwise a work will show up in the future!

	The Session structure has a flag "DevMode", with this set true all calls 
	will not generate side effects.
*/

// The credentials are used for development
// All requests using these will return
var DevCreds = Session{
	UID:     "ZeroCool",
	Session: "000-000-000",
	DevMode: true,
}

type Session struct {
	UID     string		// Your personal ID generated from theblackbox
	Session string		// The name of the current project
	DevMode bool		// If set, all uploads will be trashed, calls to finialize will not generate
				// or charge works. Keep set to true until program is functional!
				// If a work is submitted apparently unintentionally, I will not approve it
	Sessions  []string	// A list of all session tied to account

}

type Request struct {
	Action  	string
	ID      	string
	Session 	string
}

type Response struct {
	Success bool   // Status of request
	Error   string // Error for request

	ID       string   // User ID
	Session  string   // Session for action
	Sessions []string // List of all session
}

var baseURL 	= "https://theblackbox.tk/api"

func Info(auth string) (s Session, err error) {
	// Fetch new session ID
	path := fmt.Sprintf("%v?id=%v&action=info", baseURL, auth)
	res, err := fetch(path)
	if err != nil {
		return
	}

	// Check if something blew up
	if res.Success == false {
		err = errors.New(res.Error)
		return
	}

	s.UID = auth
	s.Session = res.Session
	s.Sessions = res.Sessions

	return
}

// Create a new session
// auth: Your authorization string
func Create(auth string) (s Session, err error) {
	// Fetch new session ID
	path := fmt.Sprintf("%v?id=%v&action=create", baseURL, auth)
	res, err := fetch(path)
	if err != nil {
		return
	}

	// Check if something blew up
	if res.Success == false {
		err = errors.New(res.Error)
		return
	}

	s.UID = auth
	s.Session = res.Session

	return
}

// Submit for approval
// payment: USD to be paid for work
// Finalize(100.00)
func (s *Session) Finalize(payment int) (err error) {
	// Construct and fetch response
	path := fmt.Sprintf("%v?action=finalize&pay=%v&id=%v&session=%v", baseURL, payment, s.UID, s.Session)
	if s.DevMode {
		path += "&dev=true"
	}
	res, err := fetch(path)

	if err != nil {
		return
	}

	if res.Success == false {
		err = errors.New(res.Error)
	}

	return
}

// Upload r once drained to BB
// Ensure r is NOT A STREAM
func (s *Session) Upload(r io.Reader) (err error) {
	// Construct URL
	path := fmt.Sprintf("%v?session=%v&id=%v", baseURL, s.Session, s.UID)
	if s.DevMode {
		path += "&dev=true"
	}

	// Construct multipart req using blackmagick voodoo jazz
	var b bytes.Buffer

	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", "file")
	if err != nil {
		return err
	}

	if _, err = io.Copy(fw, r); err != nil {
		return
	}

	if err = w.Close(); err != nil {
		return
	}

	req, err := http.NewRequest("POST", path, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return
	}

	// Uncan response
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var jsonRes Response
	err = json.Unmarshal(raw, &jsonRes)
	if err != nil {
		return
	}

	if jsonRes.Success != true {
		err = errors.New(jsonRes.Error)
		return
	}

	return
}

// Attach clones input from alice returning it through bob while
// feeding data to Upload. 
// This is *experimental* and simply wraps calls to Upload, in most situations
// Use of Upload is recommended.
func (s *Session) Attach(alice io.Reader) (bob io.Reader) {
	lahey := new(bytes.Buffer)
	bob = io.TeeReader(alice, lahey)

	go func() {
		s.Upload(lahey)
	}()
	return
}

// Attach clones input from alice returning it through bob while
// feeding data to Upload. 
// Attach is a convenience wrapper for Session.Attach
// This function does not require an existing Session and places
// data into the first avalible session for a given user.
func Attach(alice io.Reader, auth string) (bob io.Reader) {
	s, err := Info(auth)
	if err != nil || len(s.Sessions) == 0 {
		//Can't error, might as well give you the reader back
		return alice 
	}
	
	s.Session = s.Sessions[0]
	return s.Attach(alice)
}

func fetch(url string) (r Response, err error) {
	// Fetch and uncan
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	
	err = json.Unmarshal(raw, &r)
	return
}
