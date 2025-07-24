CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE users ADD COLUMN tenant_id UUID NOT NULL DEFAULT uuid_generate_v4();

CREATE INDEX idx_users_tenant_id ON users(tenant_id);

