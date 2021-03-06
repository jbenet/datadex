package datadex

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jbenet/data"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var webTmpl *template.Template
var webDocPages = map[string]*DocWebPage{}

const baseTmplName = "base.html"
const homeTmplName = "home.html"
const userTmplName = "user.html"
const datasetTmplName = "dataset.html"
const docTmplName = "doc.html"
const listTmplName = "list.html"

func init() {
	// templates
	webTmpl = template.New("web")
	webTmpl.Funcs(template.FuncMap{
		"timeago":   data.TimeAgo,
		"unescaped": func(s string) interface{} { return template.HTML(s) },
		"md5":       md5Hash,
		"title":     strings.Title,
		"replace": func(f, r, s string) string {
			return strings.Replace(s, f, r, -1)
		},
	})

	_, err := webTmpl.ParseGlob("web/tmpl/*.html")
	if err != nil {
		pErr("%v\n", err)
		os.Exit(-1)
	}

	// doc pages
	mdfiles, err := filepath.Glob("web/md/*.md")
	if err != nil {
		pErr("%v\n", err)
		os.Exit(-1)
	}

	for _, f := range mdfiles {
		p, err := ReadMarkdownWebPage(f)
		if err != nil {
			pErr("%v\n", err)
			os.Exit(-1)
		}

		// register page
		webDocPages[p.route] = p
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
	Packages []*Dataset
}

func webUserHandler(w http.ResponseWriter, r *http.Request) {
	u, err := requestedUser(r)
	if err != nil {
		pErr("%v\n", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	pkgs, err := indexDB.GetUserDatasets(u.Username)
	if err != nil {
		pErr("Error retrieving datasets for %s -- %v\n", u.Username, err)
		http.Error(w, "error retrieving packages", http.StatusInternalServerError)
		return
	}

	sort.Sort(sort.Reverse(DatasetsByLastUpdated(pkgs)))

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

type DatasetWebPage struct {
	I            *Dataset
	D            *data.Datafile
	Readme       string
	EscapeReadme bool
}

func webDsHomeHandler(w http.ResponseWriter, r *http.Request) {
	ds := requestDataset(r)
	h := data.NewHandle(ds)

	f, err := indexDB.GetDataset(h.Path())
	if err != nil {
		pErr("%s 404 not found\n", ds)
		http.NotFound(w, r)
		return
	}

	ref := f.Refs.ResolveRef(h.Version)
	df, err := DatafileForManifestRef(ref)
	if err != nil {
		pErr("Error loading Datafile: %s\n", err.Error())
		http.Error(w, "error loading Datafile", http.StatusInternalServerError)
		return
	}

	readme, err := FileForManifestRef(ref)
	if err != nil {
		pErr("Error finding Readme: %v\n", err)
		readme = noReadmeFile
	} else {
		if err := readme.Read(); err != nil {
			pErr("Error loading Readme: %v\n", err)
		}
	}

	webRenderPage(w, r, datasetTmplName, &WebPage{
		Title:       f.Path,
		Description: fmt.Sprintf("%s - %s", f.Path, f.Tagline),

		BodyPage: &DatasetWebPage{
			I:            f,
			D:            df,
			Readme:       string(readme.RenderedBytes()),
			EscapeReadme: !readme.CanBeRendered(),
		},
	})
}

type DocWebPage struct {
	route       string
	Title       string
	Description string
	Toc         string
	Doc         string
}

func webDocHandler(w http.ResponseWriter, r *http.Request) {
	p, found := webDocPages[r.URL.Path]
	if !found {
		pErr("%s 404 doc page not found\n", p)
		http.NotFound(w, r)
		return
	}

	webRenderPage(w, r, docTmplName, &WebPage{
		Title:       p.Title,
		Description: p.Description,
		BodyPage:    p,
	})
}

type ListWebPage struct {
	List  []interface{}
	Kind  string
	Order string
	Lists *map[string][]string
}

func webListHandler(w http.ResponseWriter, r *http.Request) {
	kind := mux.Vars(r)["kind"]
	order := mux.Vars(r)["order"]
	availableLists := &map[string][]string{
		"users":    []string{"date-registered"},
		"datasets": []string{"last-updated"},
	}

	var list []interface{}
	switch kind {
	case "users":
		users, err := indexDB.GetUsers()
		if err != nil {
			pErr("error retrieving users: %v", err)
			http.Error(w, "Error retrieving users.", http.StatusInternalServerError)
		}

		switch order {
		case "date-registered":
			sort.Sort(sort.Reverse(UsersByDateRegistered(users)))
		}

		for _, u := range users {
			list = append(list, u)
		}

	case "datasets":
		dsets, err := indexDB.GetDatasets()
		if err != nil {
			pErr("error retrieving datasets: %v", err)
			http.Error(w, "Error retrieving users.", http.StatusInternalServerError)
		}

		switch order {
		case "last-updated":
			sort.Sort(sort.Reverse(DatasetsByLastUpdated(dsets)))
			// case "downloads":
			// 	sort.Sort(DByNumDownloads(dsets))
		}

		for _, d := range dsets {
			list = append(list, d)
		}
	}

	if len(list) == 0 {
		pErr("list not found: %s/by-%s\n", kind, order)
		http.NotFound(w, r)
		return
	}

	webRenderPage(w, r, listTmplName, &WebPage{
		Title:       fmt.Sprintf("%s by %s", kind, order),
		Description: fmt.Sprintf("List of %s ordered by %s", kind, order),

		BodyPage: &ListWebPage{
			List:  list,
			Kind:  kind,
			Order: order,
			Lists: availableLists,
		},
	})
}

func webRenderPage(w http.ResponseWriter, r *http.Request,
	tmpl string, p *WebPage) {
	u := authenticatedUsername(r)
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

func webConfigHandler(w http.ResponseWriter, r *http.Request) {
	// esurl := data.ConfigGetString("db.elasticsearch", "http://localhost:9200/datadex")
	config := map[string]interface{}{
		// "ELASTICSEARCH_URL": esurl,
		"ELASTICSEARCH_URL": "/search/datadex", // use proxy pass
	}

	out, err := json.Marshal(config)
	if err != nil {
		pErr("%v\n", err)
		http.Error(w, "Server error.", http.StatusInternalServerError)
	}

	w.Write([]byte(fmt.Sprintf(webConfigScript, out)))
}

const webConfigScript = `(function() { this.datadex = this.datadex || {};
this.datadex.config = %s;
;}).call(this);
`
