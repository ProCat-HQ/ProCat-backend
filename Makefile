ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: build
build:
	GOCACHE=`pwd`/.cache go build -v ./cmd/procat

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.PHONY: dockerRun
dockerRun:
	docker run --name=procat-backend -e POSTGRES_PASSWORD=$(DB_PASSWORD) -v /tmp:/var/lib/postgresql/data -p 4321:5432 -d --rm postgres

.PHONY: dockerExec
dockerExec:
	docker exec -it $(ID) sh

.PHONY: migrationUp
migrationUp:
	migrate -path ./migrations/init -database 'postgres://postgres:$(DB_PASSWORD)@localhost:4321/postgres?sslmode=disable' up

.PHONY: migrationDown
migrationDown:
	migrate -path ./migrations/init -database 'postgres://postgres:$(DB_PASSWORD)@localhost:4321/postgres?sslmode=disable' down

.PHONY: migrationUpDownUp
migrationUpDownUp:
	make migrationUp; make migrationDown; make migrationUp

.PHONY: mockUp
mockUp:
	migrate -path ./migrations/mocks -database 'postgres://postgres:$(DB_PASSWORD)@localhost:4321/postgres?sslmode=disable' up

.PHONY: mockDown
mockDown:
	migrate -path ./migrations/mocks -database 'postgres://postgres:$(DB_PASSWORD)@localhost:4321/postgres?sslmode=disable' down

.DEFAULT_GOAL := build
