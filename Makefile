include .env

build:
	docker-compose build restapi

run:
	docker-compose up restapi

build-run:
	docker-compose up --build

stop:
	docker-compose stop

migration-up:
	 migrate -path ./migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@0.0.0.0:5432/$(POSTGRES_DB)?sslmode=disable' up

migration-down:
	migrate -path ./migrations -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@0.0.0.0:5432/$(POSTGRES_DB)?sslmode=disable' down

test:
	go test -v ./...

lint:
	golangci-lint run