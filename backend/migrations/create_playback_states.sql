CREATE TABLE playback_states (
     window_id UUID PRIMARY KEY REFERENCES windows(id) ON DELETE CASCADE,
     current_playlist_item_id UUID REFERENCES playlist_items(id) ON DELETE SET NULL,
     cycle_started_at TIMESTAMP,
     status VARCHAR(20) NOT NULL DEFAULT 'STOPPED',
     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);