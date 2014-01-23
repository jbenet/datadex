package main

import (
	"github.com/jbenet/data"
	"regexp"
)

var IndexfileNameRegexp *regexp.Regexp

func init() {
	identRE := "[A-Za-z0-9-_.]+"
	pathRE := "((" + identRE + ")/(" + identRE + "))"
	indexRE := "^" + data.DatasetDir + "/" + pathRE + "/" + IndexfileName + "$"

	IndexfileNameRegexp = compileRegexp(indexRE)
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
