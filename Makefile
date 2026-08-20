-include .env

postgres:
	docker compose up -d

postgres-down:
	docker compose down

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down -all

migratereset:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down -all
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

sqlc:
	sqlc generate

test:
	DB_SOURCE="$(DB_SOURCE)" go test -v -cover ./...

server:
	air

.PHONY: postgres postgres-down migrateup migratedown migratereset sqlc test server
