.PHONY: build
build:
	go build -o nhl-highlights -trimpath -ldflags='-s -w' main.go
