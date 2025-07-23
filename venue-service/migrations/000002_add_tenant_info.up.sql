-- add tenant_id to venues
ALTER TABLE venues ADD COLUMN tenant_id INTEGER NOT NULL DEFAULT 1;
CREATE INDEX idx_venues_tenant_id ON venues(tenant_id);

-- update unique constraint
ALTER TABLE venues DROP CONSTRAINT unique_venue_name_address;
ALTER TABLE venues ADD CONSTRAINT unique_venue_name_address_tenant 
    UNIQUE (name, address, tenant_id);
