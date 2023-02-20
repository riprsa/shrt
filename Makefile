ifneq (,$(wildcard ./.env))
	include .env
	export
endif

run:
	go run cmd/main.go
	
test:
	go test ./...