CREATE TABLE IF NOT EXISTS users (
    id bigint NOT NULL PRIMARY KEY,
    name varchar(100) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);