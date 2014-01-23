package main

import (
	"regexp"
)

func compileRegexp(s string) *regexp.Regexp {
	r, err := regexp.Compile(s)
	if err != nil {
		pOut("%s", err)
		pOut("%v", r)
		panic("Regex does not compile: " + s)
	}
	return r
}
