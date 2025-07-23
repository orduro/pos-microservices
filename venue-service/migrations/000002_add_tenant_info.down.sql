-- remove the new unique constraint
ALTER TABLE venues DROP CONSTRAINT IF EXISTS unique_venue_name_address_tenant;

-- restore the original unique constraint
ALTER TABLE venues ADD CONSTRAINT unique_venue_name_address UNIQUE (name, address);

-- remove the index
DROP INDEX IF EXISTS idx_venues_tenant_id;

-- remove tenant_id column from venues table
ALTER TABLE venues DROP COLUMN IF EXISTS tenant_id;
