ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: build
build:
	GOCACHE=`pwd`/.cache go build -v -o procat ./cmd/procat

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

# for local testing
.PHONY: dockerRun
dockerRun:
	docker run --name=procat-db -e POSTGRES_PASSWORD=$(DB_PASSWORD) -v procat_db_dev_tmp:/var/lib/postgresql/data -p $(DB_PORT):5432 -d --rm postgres

.PHONY: dockerExec
dockerExec:
	docker exec -it $(ID) sh

# for deployment
.PHONY: dockerCompose
dockerCompose:
	docker compose up --build -d

ifeq (migration,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif

# usage: make migration up 2
#		 make migration down 1
#		 make migration version
.PHONY: migration
migration:
	migrate -path ./migrations -database 'postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)' $(RUN_ARGS)

.PHONY: migrationUpDownUp
migrationUpDownUp:
	make migration up 1; make migration down 1; make migration up 1

.DEFAULT_GOAL := build
