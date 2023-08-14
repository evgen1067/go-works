COMPOSE_FILE := "./deployments/docker-compose.yaml"
COMPOSE_TEST_FILE := "./deployments/docker-compose.test.yaml"
DB_URL := "postgres://go_user:go_password@database:5432/events_db?sslmode=disable"
build:
	go build -o calendar cmd/calendar/main.go
	go build -o scheduler cmd/scheduler/main.go
	go build -o sender cmd/sender/main.go

run:
	go run cmd/calendar/main.go

lint:
	golangci-lint run ./...

test:
	go test -v -count=1 -race -timeout=30s ./...

generate:
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        api/EventService.proto

integration-tests:
	set -e ;\
    	docker compose -f ${COMPOSE_TEST_FILE} up -d --build;\
    	docker exec -it otus_go_postgres createdb --username=go_user --owner=go_user events_db ;\
    	docker compose -f ${COMPOSE_TEST_FILE} run migrate -path /migrations -database ${DB_URL} down -all ;\
        docker compose -f ${COMPOSE_TEST_FILE} run migrate -path /migrations -database ${DB_URL} up ;\
    	test_status_code=0 ;\
    	docker compose -f ${COMPOSE_TEST_FILE} run integration_tests go test || test_status_code=$$? ;\
    	docker compose -f ${COMPOSE_TEST_FILE} down ;\
    	exit $$test_status_code ;

test-cleanup:
	docker compose -f ${COMPOSE_TEST_FILE} down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker compose rm -f

up:
	docker compose -f ${COMPOSE_FILE} up

rebuild:
	docker compose -f ${COMPOSE_FILE} up --build

down:
	docker compose -f ${COMPOSE_FILE} down --remove-orphans

create_db:
	docker exec -it otus_go_postgres createdb --username=go_user --owner=go_user events_db

migrate_up:
	docker compose -f ${COMPOSE_FILE} run migrate -path /migrations -database ${DB_URL} up

migrate_down:
	docker compose -f ${COMPOSE_FILE} run migrate -path /migrations -database ${DB_URL} down

.PHONY: