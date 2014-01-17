
all: js css static ref

.PHONY: clean test

clean:
	rm -rf -- build

# css/less
css: build/static/datadex.min.css

build/static/%.min.css: build/static/%.css
	@mkdir -p `dirname $@`
	@rm -f $@
	cat $< | node_modules/.bin/cssmin > $@

build/static/%.css: other/%.less
	@mkdir -p `dirname $@`
	@rm -f $@
	node_modules/.bin/lessc $< $@

# static
static: $(shell find static -type f | sed 's/^/build\//') template.html
	cp template.html build/template.html

build/static/%: static/%
	@mkdir -p `dirname $@`
	@rm -f $@
	cp $< $@

# js/coffee
js: build/static/datadex.min.js

build/static/%.min.js: build/static/%.js
	@mkdir -p `dirname $@`
	@rm -f $@
	node_modules/.bin/uglifyjs $< > $@

build/static/datadex.js:  $(shell find js -type f | grep .js)
	@mkdir -p `dirname $@`
	@rm -f $@
	cat $^ | node_modules/.bin/uglifyjs - -b indent-level=2 -o $@

# reference
ref:
	@echo "\n\
	<!-- title: Reference -->\n\
	<!-- description: Command Reference for Data & Datadex -->\n\
	" > md/doc_ref.md
	data commands help >> md/doc_ref.md

# dependencies
deps:
	npm install
