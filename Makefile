.PHONY: clean build test start

APP_NAME = auth-api
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/db/migrations
DATABASE_URL = postgresql://auth:aaa@localhost:5432/auth?sslmode=disable
DATABASE_USER = auth
DATABASE_DB = auth

clean:
	rm -rf ./build
	go clean

build: clean
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

test:
	go test -v -cover ./...

start:
	go run ./cmd/main.go

postgres:
	docker run --name postgres13 -p 5432:5432 -e POSTGRES_USER=$(DATABASE_USER) -e POSTGRES_PASSWORD=aaa -d postgres:13

createdb:
	docker exec -it postgres13 createdb --username=$(DATABASE_USER) --owner=$(DATABASE_USER) $(DATABASE_DB)

dropdb:
	docker exec -it postgres13 dropdb --username=$(DATABASE_USER) $(DATABASE_DB)

migrate-init:
ifdef name
	migrate create -ext sql -dir $(MIGRATIONS_FOLDER) -seq ${name}
else
	@echo '"name" is not defined. Please add a name for the migration.'
endif


migrateup:
	migrate -path $(MIGRATIONS_FOLDER) -database $(DATABASE_URL) -verbose up

migrateup1:
	migrate -path $(MIGRATIONS_FOLDER) -database $(DATABASE_URL) -verbose up 1

migratedown:
	migrate -path $(MIGRATIONS_FOLDER) -database $(DATABASE_URL) -verbose down

migratedown1:
	migrate -path $(MIGRATIONS_FOLDER) -database $(DATABASE_URL) -verbose down 1
