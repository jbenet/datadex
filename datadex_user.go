package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jbenet/data"
	"io"
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

	// Auth things
	Salt     string
	PassHash string
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

	rdr, err := data.Marshal(f.Profile)
	if err != nil {
		http.Error(w, "error serializing", http.StatusInternalServerError)
		return
	}

	io.Copy(w, rdr)
}

func userAddHandler(w http.ResponseWriter, r *http.Request) {

}

func userPassHandler(w http.ResponseWriter, r *http.Request) {

}

func userAuthHandler(w http.ResponseWriter, r *http.Request) {

}

func requestUser(r *http.Request) string {
	return mux.Vars(r)["author"]
}
