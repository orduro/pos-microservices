CREATE TABLE IF NOT EXISTS venues (
   id SERIAL PRIMARY KEY,
   name TEXT NOT NULL,
   address TEXT NOT NULL,
   phone TEXT,
   venue_type TEXT,
   description TEXT NOT NULL,
   archived BOOLEAN NOT NULL DEFAULT false,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   
   CONSTRAINT unique_venue_name_address UNIQUE (name, address)
);
