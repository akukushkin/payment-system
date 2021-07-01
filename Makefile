help:
	@echo "help     - show this help"
	@echo "toolset  - setup tool for applying migration"
	@echo "setup    - apply migration with schema"
	@echo "database - start database (without schema)"
	@echo "run      - start service with database"
	@echo "e2e      - run e2e test (required running service)"
	@echo "unit     - run unit tests"

toolset:
	brew install golang-migrate

setup:
	migrate -database "postgres://payment_user:payment_pass@localhost:5432/payment_db?sslmode=disable" -path db/migrations up

database:
	docker compose up -d database

run:
	docker compose up --build

e2e:
	go test test/e2e_test.go

unit:
	go test -race -short ./...