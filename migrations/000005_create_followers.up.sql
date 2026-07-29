CREATE TABLE IF NOT EXISTS followers (
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    follower_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, follower_id),
    CHECK (user_id != follower_id)
);
