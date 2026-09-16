CREATE TABLE sync_sessions (
    id UUID PRIMARY KEY,
    media_id UUID NOT NULL REFERENCES media(id),
    started_at TIMESTAMP NOT NULL,
    duration_seconds INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);