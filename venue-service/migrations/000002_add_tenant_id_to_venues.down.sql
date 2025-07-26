-- remove the tenant-specific unique constraint
ALTER TABLE venues DROP CONSTRAINT IF EXISTS unique_venue_name_address_tenant;

-- delete duplicates of (name, address) combinations
-- delete all of them but keep the one with the
-- smallest id
DELETE FROM venues 
WHERE id NOT IN (
    SELECT MIN(id) 
    FROM venues 
    GROUP BY name, address
);

-- Now we can safely add the old unique constraint
ALTER TABLE venues ADD CONSTRAINT unique_venue_name_address UNIQUE (name, address);

-- Drop the tenant_id column
ALTER TABLE venues DROP COLUMN IF EXISTS tenant_id;
