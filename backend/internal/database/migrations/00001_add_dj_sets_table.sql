-- +goose Up
CREATE TABLE sets(
    id UUID PRIMARY KEY,
    video_id TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    dj_name TEXT NOT NULL,
    channel_name TEXT NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT Now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT Now()
);


-- +goose Down
DROP TABLE sets;
