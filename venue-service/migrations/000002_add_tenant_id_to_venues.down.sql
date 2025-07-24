ALTER TABLE venues DROP CONSTRAINT IF EXISTS unique_venue_name_address_tenant;

ALTER TABLE venues ADD CONSTRAINT unique_venue_name_address UNIQUE (name, address);

ALTER TABLE venues DROP COLUMN IF EXISTS tenant_id;
