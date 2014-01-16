package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jbenet/data"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var webTmpl *template.Template

const baseTmplName = "base.html"
const homeTmplName = "home.html"
const userTmplName = "user.html"
const datasetTmplName = "dataset.html"

func init() {
	webTmpl = template.New("web")

	// add template functions
	webTmpl.Funcs(template.FuncMap{
		"timeago": data.TimeAgo,
	})

	_, err := webTmpl.ParseGlob("web/tmpl/*.html")
	if err != nil {
		pErr("%v\n", err)
		os.Exit(-1)
	}
}

type WebPage struct {
	Title       string
	Description string
	BodyPage    interface{}

	LoggedIn bool
	Username string
}

type UserWebPage struct {
	Username string
	Name     string
	Bio      string
	Email    string
	Github   string
	Twitter  string
	Website  string
	Packages []*Indexfile
}

func webUserHandler(w http.ResponseWriter, r *http.Request) {
	u, err := requestedUserfile(r)
	if err != nil {
		pErr("%v\n", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	filenames, err := filepath.Glob(u.Dir() + "/*/" + IndexfileName)
	if err != nil {
		pErr("Error globbing indexfiles in: %s -- %v\n", u.Dir(), err)
		http.Error(w, "error retrieving packages", http.StatusInternalServerError)
		return
	}

	pkgs, err := NewIndexfiles(filenames)
	if err != nil {
		pErr("Error retrieving indexfiles in: %s -- %v\n", u.Dir(), err)
		http.Error(w, "error retrieving packages", http.StatusInternalServerError)
		return
	}

	webRenderPage(w, r, userTmplName, &WebPage{
		Title: u.User(),
		Description: fmt.Sprintf("%s has published %d datasets",
			u.User(), len(u.Profile.Packages)),

		BodyPage: &UserWebPage{
			Username: u.User(),
			Name:     u.Profile.Name,
			Email:    u.Profile.Email,
			Github:   u.Profile.Github,
			Twitter:  u.Profile.Twitter,
			Website:  u.Profile.Website,
			Packages: pkgs,
		},
	})
}

func webDsHomeHandler(w http.ResponseWriter, r *http.Request) {
	ref := mux.Vars(r)["ref"]
	url := blobUrl(ref)
	pOut("302 %v -> %v\n", ref, url)

	webRenderPage(w, r, homeTmplName, &WebPage{
		Title:       "",
		Description: "datadex.io - the dataset index",
	})
}

func webRenderPage(w http.ResponseWriter, r *http.Request,
	tmpl string, p *WebPage) {
	u := authenticatedUser(r)
	if len(u) > 0 {
		p.LoggedIn = true
		p.Username = u
	}

	err := webTmpl.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		pErr("%v\n", err)
		http.Error(w, "Error rendering page.", http.StatusInternalServerError)
	}
}
