include .env

MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	 @migrate -path=$(MIGRATIONS_PATH) -database=$(CONN_STRING) up 

.PHONY: migrate-down
migrate-down:
	 @migrate -path=$(MIGRATIONS_PATH) -database=$(CONN_STRING) down  $(filter-out $@,$(MAKECMDGOALS))
	 
.PHONY: clean-migrations
clean-migrations:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(CONN_STRING) force $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migration-status
migration-status:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(CONN_STRING) version


.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
