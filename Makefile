postgres:
	docker run --name my-pg -p 5555:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it my-pg createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it my-pg dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5555/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5555/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
