include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

test:
	echo "$${GREENLIGHT_DB_DSN}"

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


.PHONY: confirm
confirm:
	@printf "Are you sure? [y/N] "
	@read ans; [ "$${ans:-N}" = "y" ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn="$${GREENLIGHT_DB_DSN}"

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	docker exec -it postgres_greenlight_dev psql "$${GREENLIGHT_DB_DSN}"

## dev/docker-compose-up: start the development environment using docker-compose
.PHONY: dev/docker-compose-up
dev/docker-compose-up:
	docker-compose -f docker-compose.dev.yml up -d

## db/migrations/up: run all up migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo "Running up migrations..."
	migrate -path=./migrations -database="$${GREENLIGHT_DB_DSN}" up

## db/migrations/down: run the last down migration
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo "Running down migrations..."
	migrate -path=./migrations -database="$${GREENLIGHT_DB_DSN}" down 1 

## db/migrations/new: create a new migration
.PHONY: db/migrations/new
db/migrations/new: confirm
	@echo 'Creating new migration with name: $(name)'
	migrate create -seq -ext .sql -dir ./migrations ${name}

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor

## audit: run quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies'
	go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# ==================================================================================== #
# BUILD
# ==================================================================================== #


# go tool dist list: list all available GOOS and GOARCH combinations

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux/amd64/api ./cmd/api