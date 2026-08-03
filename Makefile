-include .env

postgres:
	docker compose up -d

postgres-down:
	docker compose down

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	air

.PHONY: postgres postgres-down migrateup migratedown sqlc test server
