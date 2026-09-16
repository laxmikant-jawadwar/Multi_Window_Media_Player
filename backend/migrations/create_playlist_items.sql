CREATE TABLE playlist_items (
    id UUID PRIMARY KEY,
    window_id UUID NOT NULL REFERENCES windows(id) ON DELETE CASCADE,
    media_id UUID NOT NULL REFERENCES media(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(window_id, position)
);