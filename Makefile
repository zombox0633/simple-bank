-include .env

postgres:
	docker compose up -d

postgres-down:
	docker compose down

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down -all

migratedown1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down 1

migratereset:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down -all
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

sqlc:
	sqlc generate

test:
	DB_SOURCE="$(DB_SOURCE)" go test -v -cover ./...

server:
	air

.PHONY: postgres postgres-down migrateup migrateup1 migratedown migratedown1 migratereset sqlc test server
