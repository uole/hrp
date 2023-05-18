.PHONY: build
build:
	go mod download
	CGO_ENABLED=0 go build -ldflags "-s -w" -o ./bin/hrp ./main.go