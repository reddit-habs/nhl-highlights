.PHONY: build
build:
	go build -o build/nhl-highlights -trimpath -ldflags='-s -w' main.go

.PHONY: boiler
boiler:
	rm -rf models
	sqlboiler sqlite3 --no-hooks --no-tests

install-boiler:
	go install github.com/volatiletech/sqlboiler/v4@latest
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-sqlite3@latest
