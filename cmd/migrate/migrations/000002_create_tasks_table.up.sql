CREATE TABLE IF NOT EXISTS tasks (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title varchar(255) NOT NULL,
    is_optional boolean DEFAULT FALSE,
    created_at timestamp (0) with time zone NOT NULL DEFAULT NOW()
);