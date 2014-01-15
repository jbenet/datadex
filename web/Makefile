
all: css static

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

deps:
	npm install
