# Venue Postgres connection string
# Currently using localhost for local development and db migration
# Make sure to change DSN when migrating on actual deployed db
VENUE_POSTGRES_DSN = postgres://root:toor@localhost:5432/venue_service_db?sslmode=disable
AUTH_POSTGRES_DSN = postgres://root:toor@localhost:5433/auth_service_db?sslmode=disable

# Venue database migrations
venue-db-up:
	migrate -path venue-service/migrations -database "${VENUE_POSTGRES_DSN}" -verbose up ${n}
venue-db-down:
	migrate -path venue-service/migrations -database "${VENUE_POSTGRES_DSN}" -verbose down ${n}

# Auth database migrations
auth-db-up:
	migrate -path auth-service/migrations -database "${AUTH_POSTGRES_DSN}" -verbose up ${n}
auth-db-down:
	migrate -path auth-service/migrations -database "${AUTH_POSTGRES_DSN}" -verbose down ${n}

# Migrate all databases up
db-up: venue-db-up auth-db-up

# Migrate all databases down  
db-down: venue-db-down auth-db-down

.PHONY: venue-db-up venue-db-down auth-db-up auth-db-down db-up db-down
