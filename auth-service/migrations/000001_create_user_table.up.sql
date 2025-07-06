CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    verification_token TEXT,
    verification_token_expires_at TIMESTAMPTZ,
    reset_password_token TEXT,
    reset_password_token_expires_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_verification_token ON users(verification_token);
CREATE INDEX idx_users_reset_password_token ON users(reset_password_token);
CREATE INDEX idx_users_is_verified ON users(is_verified);
