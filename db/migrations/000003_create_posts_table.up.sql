CREATE TABLE IF NOT EXISTS posts (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    description varchar(512) NOT NULL,
    -- created_at timestamp NOT NULL DEFAULT (now() AT TIME ZONE 'Asia/Jakarta'),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamp
)