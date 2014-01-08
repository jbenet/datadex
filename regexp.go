package main

import (
	"github.com/jbenet/data"
	"regexp"
)

var IndexfileNameRegexp *regexp.Regexp
var UserfileNameRegexp *regexp.Regexp

func init() {
	identRE := "[A-Za-z0-9-_.]+"
	pathRE := "((" + identRE + ")/(" + identRE + "))"
	indexRE := "^" + data.DatasetDir + "/" + pathRE + "/" + IndexfileName + "$"
	userRE := "^" + data.DatasetDir + "/" + identRE + "/" + UserfileName + "$"

	IndexfileNameRegexp = compileRegexp(indexRE)
	UserfileNameRegexp = compileRegexp(userRE)
}

func compileRegexp(s string) *regexp.Regexp {
	r, err := regexp.Compile(s)
	if err != nil {
		pOut("%s", err)
		pOut("%v", r)
		panic("Regex does not compile: " + s)
	}
	return r
}
