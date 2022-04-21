ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	go run cmd/api/main.go