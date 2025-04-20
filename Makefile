# Migration download
migrate-download:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz & mv migrate /usr/local/bin/

# Database settings
DATABASE_NAME=marcus_db
DATABASE_USER=marcus
DATABASE_PASSWORD=marcus123
DATABASE_PORT=5432
DATABASE_HOST=localhost
DATABASE_URL=postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable

MIGRATIONS_DIR=./database/migrations

# Create migrations EX: make create_migration NAME=create_atoms_table
create-migrations: 
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up 

migrate-last-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

# Restart everything
migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" drop -f
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migrate-seed-dataset:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" goto 2

# Golang runner
# Create dataset
create-from-dataset:
	@echo reading dataset files
	go run scripts/main.go 


run:
	go run cmd/main.go 

tests:
	go test -v ./...
