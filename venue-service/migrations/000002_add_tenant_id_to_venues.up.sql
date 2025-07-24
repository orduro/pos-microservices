ALTER TABLE venues ADD COLUMN tenant_id UUID;

ALTER TABLE venues DROP CONSTRAINT IF EXISTS unique_venue_name_address;

ALTER TABLE venues ADD CONSTRAINT unique_venue_name_address_tenant UNIQUE (name, address, tenant_id);
