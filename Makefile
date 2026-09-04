.PHONY: dev-api migrate migrate-down dev-web build

dev-api:
	cd backend && go run ./cmd/server

migrate:
	cd backend && go run ./cmd/migrate -seed ../random_codes_60000.json

migrate-down:
	cd backend && go run ./cmd/migrate -direction down

dev-web:
	npm run dev

build:
	npm run build
	cd backend && go build ./...
