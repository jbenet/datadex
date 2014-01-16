build: deps
	go build

deps:
	go get ./...

install: build
	go install

watch:
	-killall datadex
	make && ./datadex &
	@echo "[watching *.go for recompilation]"
	# for portability, use watchmedo -- pip install watchmedo
	@watchmedo shell-command --patterns="*.go;" --recursive --command='\
		echo; \
		date +"%Y-%m-%d %H:%M:%S"; \
		killall datadex; \
		make && ./datadex &' \
		.
