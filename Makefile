include .env

build:
	docker-compose build

run:
	docker-compose up

build-run:
	docker-compose up --build

stop:
	docker-compose stop

migration-up:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@0.0.0.0:5432/$(POSTGRES_DB)?sslmode=disable' up

migration-down:
	docker run -v $(shell pwd)/migrations:/migrations --network host migrate/migrate -path=./migrations/ -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@0.0.0.0:5432/$(POSTGRES_DB)?sslmode=disable' down -all

test:
	go test -v ./...

lint:
	golangci-lint run
