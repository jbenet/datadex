
all: js css static

.PHONY: clean test

clean:
	rm -rf -- build

# watching
WATCH=*/other/*.less;*/js/*.js;*/static/*;
watch:
	make
	@echo "[watching $(WATCH) for recompilation]"
	@# for portability, use watchmedo -- pip install watchmedo
	@watchmedo shell-command --patterns="$(WATCH)" --recursive --command='\
		echo; \
		date +"%Y-%m-%d %H:%M:%S"; \
		make' \
		.

# dependencies
deps:
	npm install


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
static: $(shell find static -type f | sed 's/^/build\//')

build/static/%: static/%
	@mkdir -p `dirname $@`
	@rm -f $@
	cp $< $@

# js/coffee
js: coffee build/static/datadex.min.js

coffee: $(shell find coffee -type f | grep .coffee | sed 's/coffee\//js-compiled\//g' | sed 's/\.coffee/\.js/g' )

js-compiled/%.js: coffee/%.coffee
	@mkdir -p `dirname $@`
	@rm -f $@
	node_modules/.bin/coffee -b -p -c $< > $@

build/static/%.min.js: build/static/%.js
	@mkdir -p `dirname $@`
	@rm -f $@
	node_modules/.bin/uglifyjs $< > $@

build/static/datadex.js:  $(shell find js -type f | grep .js) $(shell find js-compiled -type f | grep .js)
	@mkdir -p `dirname $@`
	@rm -f $@
	cat $^ | node_modules/.bin/uglifyjs - -b indent-level=2 -o $@
	cat $^ | node_modules/.bin/uglifyjs - -b indent-level=2 -o $@

# reference
ref:
	@echo "\n\
	<!-- title: Reference -->\n\
	<!-- description: Command Reference for Data & Datadex -->\n\
	" > md/doc_ref.md
	data commands help >> md/doc_ref.md
	@perl -pli -e "s/\s*$$//" md/doc_ref.md
	@printf '%s\n' "`cat md/doc_ref.md`" > md/doc_ref.md
