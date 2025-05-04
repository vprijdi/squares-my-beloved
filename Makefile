include .env

MIGRATIONS_PATH := ./cmd/migrate/migrations

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go
	
.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: gen-docs
gen-docs:
	@swag init -g ./main.go -d cmd,internal && swag fmt