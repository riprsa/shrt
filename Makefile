ifneq (,$(wildcard ./.env))
	include .env
	export
endif

run:
	go run cmd/nethttp/main.go
	
test:
	go test ./...