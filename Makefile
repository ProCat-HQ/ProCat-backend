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

# command below needs to add -v parameter
.PHONY: dockerRun
dockerRun:
	docker run --name=procat-backend -e POSTGRES_PASSWORD=$(DB_PASSWORD) -v /tmp:/var/lib/postgresql/data -p 5436:5432 -d --rm postgres

.PHONY: dockerExec
dockerExec:
	docker exec -it $(ID) sh

.PHONY: migrationUp
migrationUp:
	migrate -path ./migrations -database 'postgres://postgres:$(DB_PASSWORD)@localhost:5436/postgres?sslmode=disable' up

.PHONY: migrationDown
migrationDown:
	migrate -path ./migrations -database 'postgres://postgres:$(DB_PASSWORD)@localhost:5436/postgres?sslmode=disable' down

.PHONY: migrationUpDownUp
migrationUpDownUp:
	make migrationUp; make migrationDown; make migrationUp

.DEFAULT_GOAL := build
