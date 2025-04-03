-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username text UNIQUE not null,
    password text not null
);

-- +goose Down
DROP TABLE IF EXISTS users;

