CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username text NOT NULL UNIQUE,
    email text NOT NULL UNIQUE,
    password text NOT NULL,
    role_id bigint NOT NULL REFERENCES roles(id),
    created_at timestamp NOT NULL DEFAULT now()
);
