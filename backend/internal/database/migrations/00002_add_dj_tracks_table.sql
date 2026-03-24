-- +goose Up
CREATE TABLE tracks (
    id UUID PRIMARY KEY,
    set_id UUID NOT NULL REFERENCES sets(id) ON DELETE CASCADE,
    song_name TEXT NOT NULL,
    artist TEXT NULL,
    timestamp_in_set INTEGER NULL,
    position INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT Now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT Now()
);

-- +goose Down
DROP TABLE tracks;
