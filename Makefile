.PHONY: build
build:
	go build -o build/nhl-highlights -trimpath -ldflags='-s -w' main.go

.PHONY: boiler
boiler:
	sqlboiler sqlite3
