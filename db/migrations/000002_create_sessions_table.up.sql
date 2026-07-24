CREATE TABLE sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_email varchar(255) NOT NULL,
    refresh_token varchar(512) NOT NULL,
    is_revoked bool NOT NULL DEFAULT false,
    created_at timestamp DEFAULT now(),
    expires_at timestamp
);