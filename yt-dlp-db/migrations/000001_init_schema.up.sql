CREATE TYPE entry_status AS ENUM ('new', 'pending', 'downloading', 'completed', 'failed');

CREATE TYPE media_format AS ENUM ('mp3', 'aac', 'flac', 'wav', 'm4a', 'opus', 'vorbis', 'alac');

CREATE TABLE playlists (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    url        TEXT NOT NULL,
    title      TEXT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE playlist_entries (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    playlist_id BIGINT                   NOT NULL REFERENCES playlists (id) ON DELETE CASCADE,
    video_id    TEXT                     NOT NULL,
    title       TEXT                     NULL,
    artist      TEXT                     NULL,
    duration    INTEGER                  NULL,
    thumbnail   TEXT                     NULL,
    status      entry_status             DEFAULT 'new',
    file_path   TEXT                     NULL,
    format      media_format             DEFAULT 'mp3',
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT playlist_entries_unique UNIQUE (playlist_id, video_id)
);

CREATE INDEX idx_playlist_entries_playlist_id ON playlist_entries (playlist_id);