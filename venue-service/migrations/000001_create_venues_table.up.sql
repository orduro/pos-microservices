CREATE TABLE IF NOT EXISTS venues (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    phone TEXT,
    venue_type TEXT,
    description TEXT NOT NULL,
    archived BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_venues_name ON venues(name);
CREATE INDEX idx_venues_archived ON venues(archived);
CREATE INDEX idx_venues_venue_type ON venues(venue_type);
