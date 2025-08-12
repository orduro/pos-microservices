CREATE TABLE IF NOT EXISTS items (
    id BIGSERIAL PRIMARY KEY,
    venue_id BIGINT NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    is_available BOOLEAN NOT NULL DEFAULT true,
    image_url TEXT,
    position INTEGER NOT NULL DEFAULT 0 CHECK (position >= 0),
    modifiers JSONB NOT NULL DEFAULT '[]'::jsonb,
    archived BOOLEAN NOT NULL DEFAULT false,
    tenant_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT valid_modifiers CHECK (jsonb_typeof(modifiers) = 'array'),
    CONSTRAINT valid_price CHECK (price >= 0),
    CONSTRAINT valid_position CHECK (position >= 0),
    
    -- prevent duplicate item names within the same venue
    UNIQUE (venue_id, name, tenant_id)
);

-- indexes
CREATE INDEX idx_items_venue ON items(venue_id) WHERE archived = false;
CREATE INDEX idx_items_tenant ON items(tenant_id) WHERE archived = false;
CREATE INDEX idx_items_venue_category ON items(venue_id, category) WHERE archived = false;
CREATE INDEX idx_items_venue_position ON items(venue_id, position) WHERE archived = false;
CREATE INDEX idx_items_availability ON items(venue_id, is_available) WHERE archived = false;
CREATE INDEX idx_items_modifiers ON items USING gin (modifiers);