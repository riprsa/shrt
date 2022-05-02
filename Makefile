ifneq (,$(wildcard ./.env))
	include .env
	export
endif

echo:
	go run cmd/echo/main.go
nethttp:
	go run cmd/nethttp/main.go
test:
	go test ./...