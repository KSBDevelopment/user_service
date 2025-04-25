DB_URL=postgresql://user:password@localhost:5432/ksbd?sslmode=disable
MIGRATIONS_DIR=./migrations

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name