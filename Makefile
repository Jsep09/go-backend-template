.PHONY: run build docs sqlc migrate

## run: รัน server ใน development mode
run:
	go run ./cmd

## build: build binary
build:
	go build -o bin/app ./cmd

## docs: generate swagger docs
docs:
	swag init -g cmd/main.go -o docs
	@echo "✓ Swagger docs generated → http://localhost:3000/swagger"

## sqlc: generate sqlc queries
sqlc:
	sqlc generate

## migrate: push migration ขึ้น Supabase
migrate:
	supabase db push
