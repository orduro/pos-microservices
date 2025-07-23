-- create tenants table
CREATE TABLE tenants (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- insert default tenant for existing data
INSERT INTO tenants (name, slug) VALUES ('Default Tenant', 'default');

-- add tenant_id to users
ALTER TABLE users ADD COLUMN tenant_id INTEGER REFERENCES tenants(id) DEFAULT 1;
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
