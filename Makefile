# Venue Postgres connection string
VENUE_POSTGRES_DSN = postgres://root:toor@venue-service-postgres:5432/venue_service_db?sslmode=disable

# venue db migrate up
venue-db-up:
	migrate -path venue-service/migrations -database "${VENUE_POSTGRES_DSN}" -verbose up ${n}
	
# venue db migrate down
venue-db-down:
	migrate -path venue-service/migrations -database "${VENUE_POSTGRES_DSN}" -verbose down ${n}

.PHONY: venue-db-up venue-db-down
