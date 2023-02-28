ifneq (,$(wildcard ./.env))
	include .env
	export
endif

test:
	go test ./...

run:
	go run cmd/main.go
