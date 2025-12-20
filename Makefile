.PHONY: migrate-up migrate-down migrate-force

MIGRATION_PATH=migrations
DB_URL=postgresql://postgres:123456@localhost:5432/pentacore?sslmode=disable
MIGRATE_BIN := $(shell command -v migrate 2>/dev/null || echo "")

# Menjalankan migrasi up
migrate-up: check-migrate
	migrate -database "$(DB_URL)" -path $(MIGRATION_PATH) up

check-migrate:
	@command -v migrate >/dev/null 2>&1 || (echo "❌ migrate tidak ditemukan. Install dengan: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest" && exit 1)

# Menjalankan migrasi down (rollback)
migrate-down:
	migrate -database "$(DB_URL)" -path $(MIGRATION_PATH) down

# Force migrasi ke versi tertentu jika terjadi error
migrate-force:
	migrate -database "$(DB_URL)" -path $(MIGRATION_PATH) force $(version)

# Membuat file migrasi baru
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)

# Membuat folder domain baru dengan file entity, handler, repository, dan usecase
# Command: make domain name=users
domain:
	@if [ -z "$(name)" ]; then \
		echo "❌ Tolong isi nama domain, contoh: make domain name=users"; \
	else \
		mkdir -p internal/domain/$(name); \
		for f in entity handler repository usecase; do \
			echo "package $(name)" > internal/domain/$(name)/$$f.go; \
		done; \
		echo "✅ Domain '$(name)' berhasil dibuat di internal/domain/$(name)"; \
	fi

seed:
	psql "$(DB_URL)" -f migrations/seeders/seed_departments.sql
	psql "$(DB_URL)" -f migrations/seeders/seed_positions.sql

# Menjalankan aplikasi
run:
	go run cmd/server/main.go