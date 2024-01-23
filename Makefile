-include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## db/psql: enter a psql repl connect to database
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

## test: run all tests

.PHONY: test
test:
	go test ./... --count=1 -v
