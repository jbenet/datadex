package main

import (
	"bytes"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"path"
	"strings"
)

// Markdown stuff.

func ReadMarkdownWebPage(filename string) (*DocWebPage, error) {
	p := &DocWebPage{}
	Toc := []byte{}
	Doc := []byte{}

	// figure out route.
	p.route = "/" + path.Base(filename)
	p.route = strings.Replace(p.route, ".md", "", -1)
	p.route = strings.Replace(p.route, "_", "/", -1)

	// get the markdown first
	m, err := RenderMarkdownFile(filename)
	if err != nil {
		return nil, err
	}

	// perform surgery on the markdown.
	// this is hideous. but it will have to do.

	// split on nav, removing nav.
	{
		no := []byte("<nav>")
		nc := []byte("</nav>")
		ns := bytes.Index(m, no)
		ne := bytes.Index(m, nc)
		Toc = m[ns+len(no) : ne]
		Doc = m[ne+len(nc):]
	}

	// patch up the toc. (discard top ul)
	{
		ulo := []byte("<ul>")
		ulc := []byte("</ul>")
		uls := bytes.Index(Toc, ulo)
		ule := bytes.LastIndex(Toc, ulc)
		if uls >= 0 && ule > uls+4 {
			Toc = Toc[uls+4 : ule]
		}
	}

	// replace <ul> -> <ul class="nav">
	{
		ulo := []byte("<ul>")
		ulos := []byte("<ul class=\"nav\">")
		Toc = bytes.Replace(Toc, ulo, ulos, -1)
	}

	// figure out title
	{
		to := []byte("<!-- title: ")
		tc := []byte(" -->")
		ts := bytes.Index(Doc, to)
		if ts >= 0 {
			r := Doc[ts+len(to):]
			te := bytes.Index(r, tc)
			if te >= 0 {
				p.Title = string(r[:te])
			}
		}
	}

	// figure out description
	{
		do := []byte("<!-- description: ")
		dc := []byte(" -->")
		ds := bytes.Index(Doc, do)
		if ds >= 0 {
			r := Doc[ds+len(do):]
			de := bytes.Index(r, dc)
			if de >= 0 {
				p.Description = string(r[:de])
			}
		}
	}

	// Convert to strings
	p.Toc = string(Toc)
	p.Doc = string(Doc)

	// pOut("Parsed Markdown Page: %s\n", filename)
	// pOut("route: %s\n", p.route)
	// pOut("title: %s\n", p.Title)
	// pOut("description: %s\n", p.Description)
	// pOut("toc:\n %s\n", p.Toc)
	// pOut("doc:\n %s\n", p.Doc)

	return p, nil
}

func RenderMarkdownFile(filename string) ([]byte, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return RenderMarkdown(buf), nil
}

func RenderMarkdown(input []byte) []byte {
	// set up the HTML renderer
	flags := 0
	flags |= blackfriday.HTML_USE_XHTML
	flags |= blackfriday.HTML_USE_SMARTYPANTS
	flags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	flags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	flags |= blackfriday.HTML_SKIP_SCRIPT
	flags |= blackfriday.HTML_TOC
	// flags |= blackfriday.HTML_GITHUB_BLOCKCODE
	renderer := blackfriday.HtmlRenderer(flags, "", "")

	// set up the parser
	ext := 0
	ext |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	ext |= blackfriday.EXTENSION_TABLES
	ext |= blackfriday.EXTENSION_FENCED_CODE
	ext |= blackfriday.EXTENSION_AUTOLINK
	ext |= blackfriday.EXTENSION_STRIKETHROUGH
	ext |= blackfriday.EXTENSION_SPACE_HEADERS
	ext |= blackfriday.EXTENSION_HARD_LINE_BREAK
	ext |= blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK

	return blackfriday.Markdown(input, renderer, ext)
}
