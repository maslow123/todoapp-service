DIR = $(shell pwd)
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1 resetdb \
		infratest buildapi

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -i postgres createdb --username=root --owner=root todoapp

dropdb:
	docker exec -i postgres dropdb todoapp

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/todoapp?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/todoapp?sslmode=disable" -verbose down -all 
	
sqlc: 
	docker run --rm -v ${DIR}:/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/maslow123/todoapp-services/db/sqlc Store

infratest:
	docker-compose up -d --force-recreate testdb
	docker-compose up migrate

buildapi:
	docker-compose build --no-cache testapi

runapi:
	docker-compose up --force-recreate -d testapi

# make createsql FILENAME=
createsql:
	migrate create -ext sql -dir db/migration -seq $(FILENAME)

resetdb: dropdb createdb migrateup sqlc