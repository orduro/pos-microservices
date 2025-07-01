# load env
include .env
export

.PHONY: venue-db-up venue-db-down

# venue db migrate up
venue-db-up:
	cd venue-service && migrate -path migrations -database "${VENUE_POSTGRES_DSN}" -verbose up ${n}
	
# venue db migrate down
venue-db-down:
	cd venue-service && migrate -path migrations -database "${VENUE_POSTGRES_DSN}" -verbose down ${n}