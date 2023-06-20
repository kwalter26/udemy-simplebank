DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres-alpine --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine

redis:
	docker run --name redis -p 6379:6379 -d redis:7.0.3-alpine

createdb:
	docker exec -it postgres-alpine createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres-alpine dropdb simple_bank

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1


migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -short -v -coverpkg=./... -cover ./... -coverprofile=coverage.out -short ./tools

server:
	go run main.go

gin:
	gin -i run main.go --all --port 8080

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/kwalter26/udemy-simplebank/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	rm -f pb/*.go
	rm -f doc/swagger/api.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
        --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=api \
        proto/*.proto

evans:
	evans -r repl -p 9090

.PHONY: postgres createdb migratedown migratedown1 migrateup migrateup1 db_docs db_schema proto

