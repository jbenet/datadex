package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jbenet/data"
	"github.com/vaughan0/go-password"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

const UserfileName = "user/info"

// Userfile is the file that describes a user. It is merely a map
// of strings. Path is "datasets/<owner>/user/info"
type Userfile struct {
	data.SerializedFile "-"

	// Public profile. Viewable and settable by user.
	Profile data.UserProfile

	// Password hash, using go-password (bcrypt).
	Pass string
}

func UserfilePath(user string) string {
	return path.Join(data.DatasetDir, user, UserfileName)
}

// Constructs a new Userfile, based on its path: "<owner>/user/info"
func NewUserfile(p string) (*Userfile, error) {
	if !UserfileNameRegexp.MatchString(p) {
		return nil, fmt.Errorf("invalid Userfile path: %v", p)
	}

	f := &Userfile{SerializedFile: data.SerializedFile{Path: p}}
	f.SerializedFile.Format = f

	err := f.ReadFile()
	if err != nil {
		return f, err
	}

	return f, nil
}

func (f *Userfile) User() string {
	return strings.Split(f.Path, "/")[1]
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	u := requestUser(r)
	fmt.Fprintf(w, "%s\n", u)
}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	u := requestUser(r)

	if r.Method == "POST" {
		f, _ := NewUserfile(UserfilePath(u))

		err := data.Unmarshal(r.Body, &f.Profile)
		if err != nil {
			http.Error(w, "error serializing", http.StatusBadRequest)
			return
		}

		err = f.WriteFile()
		if err != nil {
			http.Error(w, "error saving user file", http.StatusInternalServerError)
			return
		}
	}

	f, err := NewUserfile(UserfilePath(u))
	if err != nil {
		pOut("%v\n", err)
		http.Error(w, "404 user not found", http.StatusNotFound)
		return
	}

	httpWriteMarshal(w, f.Profile)
}

func userAddHandler(w http.ResponseWriter, r *http.Request) {
	user := requestUser(r)

	email := mux.Vars(r)["email"]
	if !data.EmailRegexp.MatchString(email) {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	pass, err := ioutil.ReadAll(r.Body)
	if err != nil {
		pOut("%v\n", err)
		http.Error(w, "error with request password", http.StatusBadRequest)
		return
	}

	f, err := NewUserfile(UserfilePath(user))
	if err == nil {
		pOut("%v\n", err)
		pOut("attempt to re-register user: %s?\n", user)
		http.Error(w, "user exists", http.StatusForbidden)
		return
	}

	// ok, store user.
	f.Pass = password.Hash(string(pass[:]))
	f.Profile.Email = email

	pOut("Pass1: %s\n", pass)
	pOut("Pass2: %s\n", f.Pass)

	err = f.WriteFile()
	if err != nil {
		pOut("%v\n", err)
		http.Error(w, "error writing user file", http.StatusInternalServerError)
		return
	}

	// send verification email here...
}

func userPassHandler(w http.ResponseWriter, r *http.Request) {

}

func userAuthHandler(w http.ResponseWriter, r *http.Request) {

}

func requestUser(r *http.Request) string {
	return mux.Vars(r)["author"]
}
