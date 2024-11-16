CREATE TABLE users
(
    id                    SERIAL PRIMARY KEY,
    name                  VARCHAR(255),
    email                 VARCHAR(255) UNIQUE NOT NULL,
    uid                   VARCHAR(255) UNIQUE NOT NULL,
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at            TIMESTAMP DEFAULT NULL
);

CREATE INDEX users_uid_idx ON users (uid);
