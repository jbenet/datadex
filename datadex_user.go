package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jbenet/data"
	"github.com/vaughan0/go-password"
	"net/http"
)

// User is the object that describes a user.
type User struct {
	Username string

	// Public profile. Viewable and settable by user.
	Profile data.UserProfile

	// Password hash, using go-password (bcrypt).
	Pass string

	// Authentication token. used to verify requests.
	// (Password change clears the token.)
	AuthToken string

	// Disabled accounts cannot upload.
	Disabled bool ",omitempty"
}

func (f *User) User() string {
	return f.Username
}

func (f *User) GenerateToken() (string, error) {
	s, err := randString(20)
	if err != nil {
		return "", err
	}
	return data.StringHash(s + f.User() + f.Pass)
}

func (f *User) Put() error {
	return indexDB.PutUser(f)
}

// Route Handlers

func userHandler(w http.ResponseWriter, r *http.Request) {
	u := mux.Vars(r)["user"]
	fmt.Fprintf(w, "%s\n", u)
}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		f, err := requestedUserAuthenticated(r)
		if err != nil {
			pErr("%v\n", err)
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		err = data.Unmarshal(r.Body, &f.Profile)
		if err != nil {
			http.Error(w, "error serializing", http.StatusBadRequest)
			return
		}

		err = f.Put()
		if err != nil {
			http.Error(w, "error saving user file", http.StatusInternalServerError)
			return
		}
	}

	f, err := requestedUser(r)
	if err != nil {
		pOut("%v\n", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	httpWriteMarshal(w, f.Profile)
}

func userAddHandler(w http.ResponseWriter, r *http.Request) {
	m := &data.NewUserMsg{}
	err := data.Unmarshal(r.Body, m)
	if err != nil {
		http.Error(w, "error serializing", http.StatusInternalServerError)
		return
	}

	if len(m.Pass) < data.PasswordMinLength {
		http.Error(w, "invalid password", http.StatusBadRequest)
		return
	}

	if !data.EmailRegexp.MatchString(m.Email) {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	// not authenticated User. ignore current auth token.
	f, err := requestedUser(r)
	if err == nil {
		pOut("%v\n", err)
		pOut("attempt to re-register user: %s?\n", f.User())
		http.Error(w, "user exists", http.StatusForbidden)
		return
	}

	// ok, store user.
	f.Pass = password.Hash(string(m.Pass))
	f.Profile.Email = m.Email

	// pOut("Pass1: %s\n", m.Pass)
	// pOut("Pass2: %s\n", f.Pass)

	err = f.Put()
	if err != nil {
		pErr("%v\n", err)
		http.Error(w, "error writing user file", http.StatusInternalServerError)
		return
	}

	// send verification email here...
}

func userPassHandler(w http.ResponseWriter, r *http.Request) {
	phs := &data.NewPassMsg{}
	err := data.Unmarshal(r.Body, phs)
	if err != nil {
		http.Error(w, "error serializing", http.StatusInternalServerError)
		return
	}

	// not authenticated User. ignore current auth token.
	f, err := requestedUser(r)
	if err != nil {
		pErr("%v\n", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// pOut("Current: %s\n", phs.Current)
	// pOut("New: %s\n", phs.New)

	if !password.Check(phs.Current, f.Pass) {
		pErr("failed attempt to change password for %s\n", f.User())
		http.Error(w, "user or password incorrect", http.StatusForbidden)
		return
	}

	// ok, store new pass.
	f.Pass = password.Hash(phs.New)

	// clear AuthToken so every client needs to re-auth
	f.AuthToken = ""

	err = f.Put()
	if err != nil {
		pOut("%v\n", err)
		http.Error(w, "error writing user file", http.StatusInternalServerError)
		return
	}

	// send notification email here...
}

func userAuthHandler(w http.ResponseWriter, r *http.Request) {

	ph := ""
	err := data.Unmarshal(r.Body, &ph)
	if err != nil {
		http.Error(w, "error serializing", http.StatusInternalServerError)
		return
	}

	// not authenticated User. ignore current auth token.
	f, err := requestedUser(r)
	if err != nil {
		pErr("%v\n", err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if !password.Check(ph, f.Pass) {
		pErr("failed attempt to auth as %s\n", f.User())
		http.Error(w, "user or password incorrect", http.StatusForbidden)
		return
	}

	// Generate new token, if there is none.
	if len(f.AuthToken) == 0 {
		f.AuthToken, err = f.GenerateToken()
		if err != nil {
			pErr("Error generating token. %v\n", err)
			http.Error(w, "500 server error", http.StatusInternalServerError)
			return
		}
	}

	err = f.Put()
	if err != nil {
		pErr("Error writing user file. %v\n", err)
		http.Error(w, "500 server error", http.StatusInternalServerError)
		return
	}

	// ok, return token
	// (worry later about needing multiple tokens, etc.)
	fmt.Fprintf(w, "%s", f.AuthToken)
}

func userAwsCredHandler(w http.ResponseWriter, r *http.Request) {
	f, err := requestedUserAuthenticated(r)
	if err != nil {
		pErr("%v\n", err)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if f.Disabled {
		pErr("AwsCred request from disabled user %s forbidden.", f.User)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	c, err := getAwsFederationCredentials(f.User())
	if err != nil {
		pErr("%v\n", err)
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	pOut("user awscredentials %s\n", f.User())
	httpWriteMarshal(w, c)
}

// requested user

func requestedUser(r *http.Request) (*User, error) {
	u := mux.Vars(r)["user"]
	if len(u) == 0 {
		return nil, fmt.Errorf("No username in request.")
	}

	user, err := indexDB.GetUser(u)
	if user == nil {
		user = &User{Username: u}
	}
	return user, err
}

// request auth stuff

func requestedUserAuthenticated(r *http.Request) (*User, error) {
	f, err := authenticatedUser(r)
	if err != nil {
		return nil, err
	}

	u := mux.Vars(r)["user"]
	if u != f.User() {
		return nil, fmt.Errorf("Authenticated user is not request user.")
	}

	// ok, seems like this one's good :)
	return f, nil
}

func authenticatedUser(r *http.Request) (*User, error) {
	u := r.Header.Get(data.HttpHeaderUser)
	t := r.Header.Get(data.HttpHeaderToken)
	if len(u) < 1 || len(t) < 1 {
		return nil, fmt.Errorf("No user or token provided.")
	}

	f, err := indexDB.GetUser(u)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving User. %s", err)
	}

	// user and token must match stored values.
	if u != f.User() || t != f.AuthToken {
		return nil, fmt.Errorf("User or token mismatch.")
	}

	// ok, seems authenticated
	return f, nil
}

func authenticatedUsername(r *http.Request) string {
	f, err := authenticatedUser(r)
	if err != nil {
		return ""
	}
	return f.User()
}
