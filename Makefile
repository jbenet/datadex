build: deps
	go build
	cd datadex && go build

deps:
	go get ./...

install: build
	go install
	cd datadex && go install

WATCH=*.go;*/web/tmpl/*.html;*/web/md/*.md
watch:
	-killall datadex
	make && datadex/datadex &
	@echo "[watching $(WATCH) for recompilation]"
	# for portability, use watchmedo -- pip install watchmedo
	@watchmedo shell-command --patterns="$(WATCH)" --recursive --command='\
		echo; \
		date +"%Y-%m-%d %H:%M:%S"; \
		killall datadex; \
		make && datadex/datadex &' \
		.
