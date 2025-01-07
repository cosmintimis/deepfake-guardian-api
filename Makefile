# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./
BINARY_NAME := server.exe 
## build: build the application
.PHONY: build
build:
	go build -buildvcs=false -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run/live: run the application with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.52.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/${BINARY_NAME}" --build.delay "100" \
		--build.exclude_dir "frontend" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"