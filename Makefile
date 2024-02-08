-include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run/inventory
.PHONY: run/inventory
run/inventory:
	go mod tidy
	go run ./cmd/inventory

## run/automigrate
.PHONY: run/automigrate
run/automigrate:
	go run ./cmd/automigrate --dsn ${DB_DSN}

## db/psql: enter a psql repl connect to database
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

## test: run all tests
.PHONY: test
test:
	go test ./... --count=1 -v

## gen-mock: generate mocks
.PHONY: gen-mock
gen-mock:
	cd order; mockery
