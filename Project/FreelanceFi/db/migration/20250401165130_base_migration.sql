-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    role TEXT NOT NULL DEFAULT 'freelancer'
    
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE sessions;
DROP TABLE users;
ALTER TABLE users ADD COLUMN password TEXT NOT NULL;
ALTER TABLE users ADD COLUMN password_hash VARCHAR(255);
ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'freelancer';


