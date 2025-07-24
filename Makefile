APP_NAME=storage
PATH_TO_MAIN=./cmd/storage/main.go
PATH_TO_CONFIG=./config/local.yaml
PATH_TO_MIGRATOR=./cmd/migrator/main.go
MIGRATIONS_TEST_TABLE=migrations_test

.PHONY: build run migrator migrator_test test clean

build:
	go build -o ./bin/$(APP_NAME) $(PATH_TO_MAIN) -config=$(PATH_TO_CONFIG)

run:
	go run $(PATH_TO_MAIN) -config $(PATH_TO_CONFIG)

migrator:
	go run $(PATH_TO_MIGRATOR) -config $(PATH_TO_CONFIG)

migrator_test:
	go run $(PATH_TO_MIGRATOR) -config=$(PATH_TO_CONFIG) --migrations-table=$(MIGRATIONS_TEST_TABLE)

test:
	go test ./...

clean:
	rm -rf bin/