postgres:
	docker run --name postgres-alpine -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine
createdb:
	docker exec -it postgres-alpine createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres-alpine dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...  -coverprofile=coverage.out

server:
	go run main.go

gin:
	gin -i run main.go --all --port 8080

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/kwalter26/udemy-simplebank/db/sqlc Store

.PHONY: postgres createdb

