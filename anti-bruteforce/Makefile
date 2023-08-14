build:
	go build -o bruteforce cmd/bruteforce/main.go

run:
	go run cmd/bruteforce/main.go

lint:
	golangci-lint run ./...

test:
	go test -v -count=1 -race -timeout=30s ./...

up:
	docker compose -f ./deployments/docker-compose.yaml up

rebuild:
	docker compose -f ./deployments/docker-compose.yaml up --build

down:
	docker compose -f ./deployments/docker-compose.yaml down --remove-orphans

create_db:
	docker exec -it otus_go_postgres createdb --username=go_user --owner=go_user bruteforce_db

migrate_up:
	docker compose -f ./deployments/docker-compose.yaml run migrate -path /migrations -database "postgres://go_user:go_password@database:5432/bruteforce_db?sslmode=disable" up

migrate_down:
	docker compose -f ./deployments/docker-compose.yaml run migrate -path /migrations -database "postgres://go_user:go_password@database:5432/bruteforce_db?sslmode=disable" down

.PHONY: