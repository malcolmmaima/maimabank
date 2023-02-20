postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=7915 -d postgres:12-alpine

createdb: 
	docker exec -it postgres12 createdb --username=root --owner=root maima_bank

dropdb: 
	docker exec -it postgres12 dropdb maima_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:7915@localhost:5432/maima_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:7915@localhost:5432/maima_bank?sslmode=disable" -verbose down

sqlc: 
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/Store.go github.com/malcolmmaima/maimabank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock