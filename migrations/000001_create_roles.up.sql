CREATE TABLE IF NOT EXISTS roles (
    id bigserial PRIMARY KEY,
    name text NOT NULL UNIQUE
);
